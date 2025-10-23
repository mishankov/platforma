package session

import (
	"context"
	"database/sql"

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

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return []database.Migration{}, database.Schema{
		Queries: []string{
			`
		CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(255) PRIMARY KEY,
			"user" VARCHAR(255),
			created TIMESTAMP,
			expires TIMESTAMP
		)
		`,
		},
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*Session, error) {
	var session Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) GetByUserId(ctx context.Context, userID string) (*Session, error) {
	var session Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) Create(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (id, "user", created, expires)
		VALUES (:id, :user, :created, :expires)
	`
	_, err := r.db.NamedExecContext(ctx, query, session)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
