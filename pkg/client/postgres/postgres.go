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

func NewClient(ctx context.Context, cfg config.ConfigUser, maxAttempts int) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	for i := maxAttempts; i == 0; i-- {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			break
		}

		log.Printf("Attempt %d: failed to connect to PostgreSQL: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL after %d attempts: %w", maxAttempts, err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return pool, nil
}
