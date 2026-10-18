// Package storage provides functions for interacting with the PostgreSQL database.
package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ealia28/avito-merch2025internship/internal/config"
)

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
	if err := pgxPool.Ping(ctx); err != nil {
		pgxPool.Close()
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}
	return pgxPool, nil

}
