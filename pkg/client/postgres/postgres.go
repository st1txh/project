package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"rest-api-tutorial/internal/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arg ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arg ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, cfg config.User, maxAttempts int) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxAttempts; i++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		pool, err = pgxpool.Connect(attemptCtx, dsn)
		cancel() // Важно освобождать ресурсы контекста

		if err == nil {
			// Проверяем работоспособность подключения
			conn, err := pool.Acquire(ctx)
			if err != nil {
				log.Printf("Connection acquired but failed to ping: %v", err)
				continue
			}
			conn.Release()

			log.Println("Successfully connected to PostgreSQL!")
			return pool, nil
		}

		log.Printf("Attempt %d/%d failed: %v", i+1, maxAttempts, err)
		time.Sleep(time.Second * time.Duration(i+1)) // Увеличиваем задержку с каждой попыткой
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, err)
}
