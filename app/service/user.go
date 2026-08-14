package service

import (
	"fmt"
	"errors"
	"context"

	"github.com/nekto-sns/nekto-server/app/model"
)

type userRepository interface{
	ByName(context.Context, string) (*model.User, error)
}

type userService struct{
	repo userRepository
}

func NewUser(repo userRepository) (*userService) {
	return &userService{
		repo: repo,
	}
}

func (u *userService) ByName(ctx context.Context, name string) (*model.User, error) {
	user, err := u.repo.ByName(ctx, name)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("User not found: %w", err)
		}

		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	return user, nil
}
