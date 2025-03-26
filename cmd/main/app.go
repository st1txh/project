package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api-tutorial/internal/config"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/pkg/client/postgres"
	"rest-api-tutorial/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	pgCfg := config.ConfigUser{
		Host:     cfg.PostgreSQL.Host,
		Port:     cfg.PostgreSQL.Port,
		Username: cfg.PostgreSQL.Username,
		Password: cfg.PostgreSQL.Password,
		Database: cfg.PostgreSQL.Database,
	}

	pool, err := postgres.NewClient(context.Background(), pgCfg, 5)
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()

	userStorage := user.NewUserStorage(pool, logger)
	userHandler := user.NewHandler(userStorage, logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Сервер запущен!"})
	})

	api := router.Group("/api")
	{
		api.GET("/users", userHandler.GetList)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:uuid", userHandler.UpdateUser)
		api.PATCH("/users/:uuid", userHandler.PartiallyUpdateUser)
		api.DELETE("/users/:uuid", userHandler.DeleteUser)
	}
	fmt.Printf("Api %s", api)
	start(router, cfg)
}

func start(router *gin.Engine, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Starting application...")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("Detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("Listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
		listener, listenErr = net.Listen("tcp", address)
		logger.Infof("Server is listening port %s", address)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
