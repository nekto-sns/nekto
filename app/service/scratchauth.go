package service

import (
	"fmt"
	"errors"
	"context"

	"github.com/nekto-sns/nekto-server/app/model"
)

type scratchAuthRepository interface{
	UserIDByScratchName(context.Context, string) (string, error)
}

type scratchAuthClient interface{
	Verify(context.Context, string) (string, error)
}

type sessionManager interface{
	Create(context.Context, string) (string, error)
	MaxAge() int
}

type scratchAuthService struct{
	repo      scratchAuthRepository
	saClient  scratchAuthClient
	sessionMN sessionManager
}

func NewScratchAuth(saRepo scratchAuthRepository, saClient scratchAuthClient, sessionMN sessionManager) (*scratchAuthService) {
	return &scratchAuthService{
		repo: saRepo,
		sessionMN: sessionMN,
		saClient: saClient,
	}
}

func (s *scratchAuthService) LoginCallback(ctx context.Context, privateCode string) (string, error) {
	scratchName, err := s.saClient.Verify(ctx, privateCode)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", fmt.Errorf("Scratch account not found: %w", err)
		}
		return "", fmt.Errorf("Authentication failed: %w", err)
	}

	userID, err := s.repo.UserIDByScratchName(ctx, scratchName)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", fmt.Errorf("User not found: %w", err)
		}
		return "", fmt.Errorf("Failed to get user: %w", err)
	}

	sessionID, err := s.sessionMN.Create(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("Failed to create session: %w", err)
	}

	return sessionID, err
}
