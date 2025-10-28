package session

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mishankov/platforma/database"
)

type db interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Repository struct {
	db db
}

func NewRepository(db db) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrations() []database.Migration {
	return []database.Migration{{
		ID: "init",
		Up: `CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(255) PRIMARY KEY,
			"user" VARCHAR(255),
			created TIMESTAMP,
			expires TIMESTAMP
		)`,
		Down: "DROP TABLE sessions",
	}}
}

func (r *Repository) Get(ctx context.Context, id string) (*Session, error) {
	var session Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}
	return &session, nil
}

func (r *Repository) GetByUserId(ctx context.Context, userID string) (*Session, error) {
	var session Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by user id: %w", err)
	}
	return &session, nil
}

func (r *Repository) Create(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (id, "user", created, expires)
		VALUES (:id, :user, :created, :expires)
	`
	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
