package repository

import (
	"fmt"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nekto-sns/nekto-server/app/model"
)

type authRepository struct{
	db *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) (*authRepository) {
	return &authRepository{
		db: pool,
	}
}

func (r *authRepository) UsernameByScratchName(ctx context.Context, scratchName string) (string, err) {
	var username string

	err := pool.QueryRow(ctx,
			     `SELECT u.name FROM users u
			     INNER JOIN scratch_auth sa ON u.id = sa.user_id
			     WHERE sa.scratch_name = $1`,
			scratchName,
	).Scan(&username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("Row not found: %w", model.ErrNotFound)
		}
		return "", fmt.Errorf("DB query failed (%v): %w", err, model.ErrInternal)
	}
}

