package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) migrations() []Migration {
	return []Migration{{
		ID:   "init",
		Up:   "CREATE TABLE IF NOT EXISTS platforma_migrations (repository TEXT, id TEXT, timestamp TIMESTAMP)",
		Down: "DROP TABLE platforma_migrations",
	}}
}

func (r *repository) getMigrationLogs(ctx context.Context) ([]migrationLog, error) {
	var migrations []migrationLog
	err := r.db.SelectContext(ctx, &migrations, "SELECT * FROM platforma_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get migration logs: %w", err)
	}

	return migrations, nil
}

func (r *repository) saveMigrationLog(ctx context.Context, log migrationLog) error {
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

func (r *repository) executeQuery(ctx context.Context, query string) error {
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
