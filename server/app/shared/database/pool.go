package database

import (
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create DB pool: %v\n", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to DB: %v", err)
	}
	return pool, nil
}
