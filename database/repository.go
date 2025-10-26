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

func (r *repository) GetMigrations(ctx context.Context) ([]*migrations, error) {
	var migrations []*migrations
	err := r.db.SelectContext(ctx, &migrations, "SELECT * FROM platforma_migrations")
	if err != nil {
		return nil, err
	}

	return migrations, nil
}

func (r *repository) SaveMigration(ctx context.Context, migration migrations) error {
	query := `
		INSERT INTO platforma_migrations (repository, id, timestamp)
		VALUES (:repository, :id, :timestamp)
	`
	_, err := r.db.NamedExecContext(ctx, query, migration)
	return err
}
