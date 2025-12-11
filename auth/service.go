package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/platforma-dev/platforma/log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, id, password, salt string) error
	Delete(ctx context.Context, id string) error
}

type authStorage interface {
	GetUserIdFromSessionId(context.Context, string) (string, error)
	CreateSessionForUser(context.Context, string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
	DeleteSessionsByUserId(ctx context.Context, userId string) error
}

type Service struct {
	repo              repository
	authStorage       authStorage
	sessionCookieName string
	usernameValidator func(string) error
	passwordValidator func(string) error
}

func NewService(repo repository, authStorage authStorage, sessionCookieName string, usernameValidator, passwordValidator func(string) error) *Service {
	if usernameValidator == nil {
		usernameValidator = defaultUsernameValidator
	}

	if passwordValidator == nil {
		passwordValidator = defaultPasswordValidator
	}

	return &Service{
		repo:              repo,
		authStorage:       authStorage,
		sessionCookieName: sessionCookieName,
		usernameValidator: usernameValidator,
		passwordValidator: passwordValidator,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *Service) GetFromSession(ctx context.Context, sessionId string) (*User, error) {
	userId, err := s.authStorage.GetUserIdFromSessionId(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id from request: %w", err)
	}

	if userId == "" {
		return nil, ErrUserNotFound
	}

	return s.Get(ctx, userId)
}

func (s *Service) CreateWithLoginAndPassword(ctx context.Context, username, password string) error {
	err := s.usernameValidator(username)
	if err != nil {
		return errors.Join(ErrInvalidUsername, err)
	}

	err = s.passwordValidator(password)
	if err != nil {
		return errors.Join(ErrInvalidPassword, err)
	}

	salt := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+":"+salt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	user := &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashedPassword),
		Salt:     salt,
		Created:  time.Now(),
		Updated:  time.Now(),
		Status:   StatusActive,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Service) CreateSessionFromUsernameAndPassword(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", ErrWrongUserOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+":"+user.Salt))
	if err != nil {
		return "", ErrWrongUserOrPassword
	}

	session, err := s.authStorage.CreateSessionForUser(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return session, nil
}

func (s *Service) DeleteSession(ctx context.Context, sessionId string) error {
	err := s.authStorage.DeleteSession(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *Service) CookieName() string {
	return s.sessionCookieName
}

func (s *Service) ChangePassword(ctx context.Context, currentPassword, newPassword string) error {
	user := UserFromContext(ctx)
	if user == nil {
		return ErrUserNotFound
	}

	if s.passwordValidator != nil {
		log.DebugContext(ctx, "validating new password")
		err := s.passwordValidator(newPassword)
		if err != nil {
			return errors.Join(ErrInvalidPassword, err)
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword+":"+user.Salt))
	log.DebugContext(ctx, "password validation results", "error", err)
	if err != nil {
		return ErrCurrentPasswordIncorrect
	}

	newSalt := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword+":"+newSalt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	err = s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword), newSalt)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context) error {
	user := UserFromContext(ctx)
	if user == nil {
		return ErrUserNotFound
	}

	// Delete all user sessions first
	err := s.authStorage.DeleteSessionsByUserId(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	// Delete the user
	err = s.repo.Delete(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func defaultPasswordValidator(password string) error {
	if len(password) < 8 {
		return ErrShortPassword
	}

	if len(password) > 100 {
		return ErrLongPassword
	}

	return nil
}

func defaultUsernameValidator(username string) error {
	if len(username) < 5 {
		return ErrShortUsername
	}

	if len(username) > 20 {
		return ErrLongUsername
	}

	return nil
}
