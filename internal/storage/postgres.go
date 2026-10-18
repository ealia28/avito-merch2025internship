// Package storage provides functions for interacting with the PostgreSQL database.
package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ealia28/avito-merch2025internship/internal/config"
)

const delay = 100 * time.Millisecond
const maxDelay = 2 * time.Second
const maxAttempts = 10

// NewPool creates a new pgxpool.Pool instance for the given database configuration.
// It ensures that the pool is properly initialized and can be used for database operations.
func NewPool(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres DSN: %w", err)
	}
	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	initialDelay := delay
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = pgxPool.Ping(ctx)
		if err == nil {
			return pgxPool, nil
		}
		fmt.Printf("attempt number %d failed: %v\n", attempt, err)

		if attempt == maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			pgxPool.Close()
			return nil, fmt.Errorf("ping pgx pool: context canceled %w", ctx.Err())
		case <-time.After(initialDelay):
			initialDelay *= 2
			if initialDelay > maxDelay {
				initialDelay = maxDelay
			}
		}
	}
	pgxPool.Close()
	return nil, fmt.Errorf("ping postgres after %d attempts: %w", maxAttempts, err)
}
