package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Schema() ([]Migration, Schema) {
	return []Migration{}, Schema{Queries: []string{
		"CREATE TABLE IF NOT EXISTS platforma_migrations (repository TEXT, id TEXT, timestamp TIMESTAMP)",
	}}
}

func (r *repository) GetMigrationLogs(ctx context.Context) ([]migrationLog, error) {
	var migrations []migrationLog
	err := r.db.SelectContext(ctx, &migrations, "SELECT * FROM platforma_migrations")
	if err != nil {
		return nil, err
	}

	return migrations, nil
}

func (r *repository) SaveMigrationLog(ctx context.Context, log migrationLog) error {
	query := `
		INSERT INTO platforma_migrations (repository, id, timestamp)
		VALUES (:repository, :id, :timestamp)
	`
	_, err := r.db.NamedExecContext(ctx, query, log)
	return err
}

func (r *repository) RemoveMigrationLog(ctx context.Context, repository, id string) error {
	query := `DELETE FROM platforma_migrations WHERE repository = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, repository, id)
	return err
}
