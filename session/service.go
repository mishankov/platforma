package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, session *Session) error {
	return s.repo.Create(ctx, session)
}

func (s *Service) Get(ctx context.Context, id string) (*Session, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByUserId(ctx context.Context, id string) (*Session, error) {
	return s.repo.GetByUserId(ctx, id)
}

func (s *Service) DeleteSession(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetUserIdFromSessionId(ctx context.Context, id string) (string, error) {
	session, err := s.repo.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	return session.User, nil
}

func (s *Service) CreateSessionForUser(ctx context.Context, userId string) (string, error) {
	session := &Session{
		ID:      uuid.NewString(),
		User:    userId,
		Created: time.Now(),
		Expires: time.Now().Add(100 * 24 * time.Hour),
	}

	err := s.repo.Create(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return session.ID, nil
}

func (s *Service) DeleteSessionsByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteByUserId(ctx, userId)
}
