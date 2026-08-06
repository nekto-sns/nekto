package service

import (
	"fmt"
	"errors"
	"context"

	"github.com/nekto-sns/nekto-server/app/model"
)

type authRepository interface{
	UserIDByScratchName(ctx context.Context, scratchName string) (string, error)
}

type scratchAuthClient interface{
	Verify(ctx context.Context, privateCode string) (string, bool, error)
}

type authService struct{
	repository authRepository
	scratchAuth scratchAuthClient
}

func NewAuthService(repo authRepository, saClient scratchAuthClient) (*authService) {
	return &authService{
		repository: repo,
		scratchAuth: saClient,
	}
}

func (s *authService) LoginCallback(ctx context.Context, privateCode string) (string, bool, error) {
	scratchName, isValid, err := s.scratchAuth.Verify(ctx, privateCode)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", false, fmt.Errorf("Scratch account not found: %w", err)
		}
		return "", false, fmt.Errorf("Authentication failed: %w", err)
	}

	userID, err := s.repository.UserIDByScratchName(ctx, scratchName)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", false, fmt.Errorf("User not found: %w", err)
		}
		return "", false, fmt.Errorf("Failed to get user: %w", err)
	}

	return userID, isValid, nil
}
