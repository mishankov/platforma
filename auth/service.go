package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mishankov/platforma/log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, id, password, salt string) error
}

type authStorage interface {
	GetUserIdFromSessionId(context.Context, string) (string, error)
	CreateSessionForUser(context.Context, string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
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
	return s.repo.Get(ctx, id)
}

func (s *Service) GetFromSession(ctx context.Context, sessionId string) (*User, error) {
	userId, err := s.authStorage.GetUserIdFromSessionId(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id from request: %w", err)
	}

	if userId == "" {
		return nil, nil
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
		return err
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

	return s.repo.Create(ctx, user)
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
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	return session, err
}

func (s *Service) DeleteSession(ctx context.Context, sessionId string) error {
	return s.authStorage.DeleteSession(ctx, sessionId)
}

func (s *Service) CookieName() string {
	return s.sessionCookieName
}

func (s *Service) ChangePassword(ctx context.Context, currentPassword, newPassword string) error {
	user := UserFromContext(ctx)
	if user == nil {
		return errors.New("user not found")
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
		return err
	}

	return s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword), newSalt)
}

func defaultPasswordValidator(password string) error {
	if len(password) < 8 {
		return errors.New("short password")
	}

	if len(password) > 100 {
		return errors.New("long password")
	}

	return nil
}

func defaultUsernameValidator(username string) error {
	if len(username) < 5 {
		return errors.New("short username")
	}

	if len(username) > 20 {
		return errors.New("long username")
	}

	return nil
}
