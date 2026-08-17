package service

import (
	"context"
	"crypto/rand"
	"time"
	"fmt"
	"encoding/base64"
	"errors"

	"github.com/nekto-sns/nekto-server/app/model"
)

type sessionRepository interface{
	Save(context.Context, string, string, time.Duration) error
	UserID(context.Context, string) (string, error)
}

type sessionService struct{
	repo sessionRepository
	ttl  time.Duration
}

func NewSession(repo sessionRepository, ttl time.Duration) (*sessionService) {
	return &sessionService{
		repo: repo,
		ttl: ttl,
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *sessionService) Create(ctx context.Context, userID string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", fmt.Errorf("Failed to generate session ID: %w", model.ErrInternal)
	}

	err = s.repo.Save(ctx, sessionID, userID, s.ttl)
	if err != nil {
		return "", fmt.Errorf("Failed to save session ID: %w", err)
	}

	return sessionID, nil
}

func (s *sessionService) UserID(ctx context.Context, sessionID string) (string, error) {
	sessionID, err := s.repo.UserID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", fmt.Errorf("Session not found: %w", err)
		}
		return "", fmt.Errorf("Failed to get session: %w", err)
	}

	return sessionID, nil
}
