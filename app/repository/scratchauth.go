package repository

import (
	"fmt"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nekto-sns/nekto-server/app/model"
)

type scratchAuthRepository struct{
	db *pgxpool.Pool
}

func NewScratchAuth(pool *pgxpool.Pool) (*scratchAuthRepository) {
	return &scratchAuthRepository{
		db: pool,
	}
}

func (r *scratchAuthRepository) UserIDByScratchName(ctx context.Context, scratchName string) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx,`SELECT user_id FROM scratch_auth WHERE scratch_name = $1`, scratchName).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("Row not found: %w", model.ErrNotFound)
		}
		return "", fmt.Errorf("DB query failed (%v): %w", err, model.ErrInternal)
	}

	return userID, nil
}

