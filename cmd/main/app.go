package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api-tutorial/internal/config"
	"rest-api-tutorial/internal/films"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/pkg/client/postgres"
	"rest-api-tutorial/pkg/logging"
	"time"
)

// @title Movie REST API
// @version 1.0
// @description API для управления фильмами и пользователями с аутентификацией

// @contact.name API Support
// @contact.email st1txh.devops@ussr.com

// @license.name Omsk
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in query
// @name Authorization

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	// Создаем контекст с таймаутом для инициализации приложения
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Настройка подключения к PostgreSQL
	pgCfg := config.User{
		Host:     cfg.PostgreSQL.Host,
		Port:     cfg.PostgreSQL.Port,
		Username: cfg.PostgreSQL.Username,
		Password: cfg.PostgreSQL.Password,
		Database: cfg.PostgreSQL.Database,
	}

	// Подключаемся к PostgreSQL с повторными попытками
	pool, err := postgres.NewClient(initCtx, pgCfg, 5)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Проверяем соединение
	if err := checkPostgreSQLConnection(initCtx, pool); err != nil {
		logger.Fatalf("PostgreSQL connection check failed: %v", err)
	}

	// Инициализация слоев приложения
	userStorage := user.NewUserStorage(pool, logger)
	userHandler := user.NewHandler(userStorage, logger)

	filmStorage := films.NewFilmStorage(pool, logger)
	filmHandler := films.NewHandler(filmStorage, logger)
	// Настройка роутера
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Настройка Swagger
	router.StaticFile("/swagger.json", "./docs/swagger.json") // Явное указание JSON-файла
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.json"), // Указываем прямой путь к файлу
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))
	// Middleware для добавления пула соединений в контекст
	router.Use(func(c *gin.Context) {
		c.Set("postgres_pool", pool)
		c.Next()
	})

	// Маршруты
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Сервер запущен!", "status": "ok"})
	})

	api := router.Group("/api")
	{
		api.GET("/users", userHandler.GetList)
		api.GET("/users/:uuid", userHandler.GetUser)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:uuid", userHandler.UpdateUser)
		api.PATCH("/users/:uuid", userHandler.PartiallyUpdateUser)
		api.DELETE("/users/:uuid", userHandler.DeleteUser)

		api.POST("/films", filmHandler.CreateFilm)
		api.GET("/films", filmHandler.GetList)
		api.GET("/films/sort", filmHandler.GetListSort)
		api.GET("/films/:uuid", filmHandler.GetUserFilm)
		api.PATCH("/films/:uuid", filmHandler.PartiallyUpdateFilm)
		api.DELETE("/films/:uuid", filmHandler.DeleteFilm)
	}

	// Запуск сервера
	if err := startServer(router, cfg, logger); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func checkPostgreSQLConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	if err := conn.Conn().Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

func startServer(router *gin.Engine, cfg *config.Config, logger *logging.Logger) error {
	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return fmt.Errorf("failed to get app directory: %w", err)
		}

		socketPath := path.Join(appDir, "app.sock")

		// Удаляем старый сокет, если существует
		if _, err := os.Stat(socketPath); err == nil {
			if err := os.Remove(socketPath); err != nil {
				return fmt.Errorf("failed to remove old socket: %w", err)
			}
		}

		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
		listener, listenErr = net.Listen("tcp", address)
		logger.Infof("Server is listening on %s", address)
	}

	if listenErr != nil {
		return fmt.Errorf("failed to create listener: %w", listenErr)
	}

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Application started successfully")
	return server.Serve(listener)
}
