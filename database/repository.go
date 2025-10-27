package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrations() []Migration {
	return []Migration{{
		ID:   "init",
		Up:   "CREATE TABLE IF NOT EXISTS platforma_migrations (repository TEXT, id TEXT, timestamp TIMESTAMP)",
		Down: "DROP TABLE platforma_migrations",
	}}
}

func (r *Repository) GetMigrationLogs(ctx context.Context) ([]MigrationLog, error) {
	var migrations []MigrationLog
	err := r.db.SelectContext(ctx, &migrations, "SELECT * FROM platforma_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get migration logs: %w", err)
	}

	return migrations, nil
}

func (r *Repository) SaveMigrationLog(ctx context.Context, log MigrationLog) error {
	query := `
		INSERT INTO platforma_migrations (repository, id, timestamp)
		VALUES (:repository, :id, :timestamp)
	`
	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to save migration log: %w", err)
	}
	return nil
}

func (r *Repository) RemoveMigrationLog(ctx context.Context, repository, id string) error {
	query := `DELETE FROM platforma_migrations WHERE repository = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, repository, id)
	if err != nil {
		return fmt.Errorf("failed to remove migration log: %w", err)
	}
	return nil
}

func (r *Repository) ExecuteQuery(ctx context.Context, query string) error {
	_, err := r.db.ExecContext(ctx, query)
	return err
}
