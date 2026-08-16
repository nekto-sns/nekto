package repository

import (
	"fmt"
	"errors"
	"time"
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/nekto-sns/nekto-server/app/model"
)

type sessionRepository struct{
	rdb *redis.Client
}

func NewSession(rdb *redis.Client) (*sessionRepository) {
	return &sessionRepository{
		rdb: rdb,
	}
}

func (r *sessionRepository) Save(ctx context.Context, sessionID, userID string, ttl time.Duration) error {
	err := r.rdb.Set(ctx, sessionID, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("Redis operation failed (%v): %w", err, model.ErrInternal)
	}

	return nil
}

func (r *sessionRepository) UserID(ctx context.Context, sessionID string) (string, error) {
	userID, err := r.rdb.Get(ctx, sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("Key not found: %w", model.ErrNotFound)
		}
		return "", fmt.Errorf("Redis operation failed (%v): %w", err, model.ErrInternal)
	}

	return userID, nil
}
