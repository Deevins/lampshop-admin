package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool создаёт pgx-пулл подключений к PostgreSQL по URL в env PG_URL.
func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig("postgres://postgres:secret@localhost:5436/lampshop_admin?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	// Можно настроить пул (максимум соединений, таймауты и т.п.)
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create pgxpool: %w", err)
	}

	// Убедимся, что соединение работает:
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
