package sessionmanager

import (
	"fmt"
	"context"
	"time"
	"crypto/rand"
	"encoding/base64"

	"github.com/nekto-sns/nekto/server/app/model"
)

type sessionRepository interface{
	Save(context.Context, string, string, time.Duration) error
}

type sessionManager struct{
	repo   sessionRepository
	ttl    time.Duration
	maxAge int
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func New(repo sessionRepository, ttl time.Duration) *sessionManager {

	return &sessionManager{
		repo: repo,
		ttl:  ttl,
		maxAge: int(ttl.Seconds()),
	}
}

func (m *sessionManager) Create(ctx context.Context, userID string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", fmt.Errorf("Failed to generate session (%v): %w", err, model.ErrInternal)
	}

	err = m.repo.Save(ctx, sessionID, userID, m.ttl)
	if err != nil {
		return "", fmt.Errorf("Failed to save session ID: %w", err)
	}

	return sessionID, nil
}

func (m *sessionManager) MaxAge() int {
	return m.maxAge
}
