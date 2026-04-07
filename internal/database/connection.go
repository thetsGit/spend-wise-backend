package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetsGit/spend-wise-be/internal/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(config *config.Config) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), config.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Ping to verify connection status
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close shuts down the pool, can call this when server stops
func (db *DB) Close() {
	db.Pool.Close()
}
