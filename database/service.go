package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type service struct {
	repo *repository
	db   *sqlx.DB
}

func newService(repo *repository, db *sqlx.DB) *service {
	return &service{repo: repo, db: db}
}

func (s *service) GetMigrationLogs(ctx context.Context) ([]migrationLog, error) {
	return s.repo.GetMigrationLogs(ctx)
}

func (s *service) SaveMigrationLog(ctx context.Context, migration migrationLog) error {
	return s.repo.SaveMigrationLog(ctx, migration)
}

func (s *service) RemoveMigrationLog(ctx context.Context, repository, id string) error {
	return s.repo.RemoveMigrationLog(ctx, repository, id)
}

func (s *service) ApplyMigration(ctx context.Context, migration Migration) error {
	_, err := s.db.ExecContext(ctx, migration.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}
	return nil
}

func (s *service) RevertMigration(ctx context.Context, migration Migration) error {
	_, err := s.db.ExecContext(ctx, migration.Down)
	if err != nil {
		return fmt.Errorf("failed to revert migration: %w", err)
	}
	return nil
}

func (s *service) ApplySchema(ctx context.Context, schema Schema) error {
	for _, query := range schema.Queries {
		_, err := s.db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to apply schema query: %w", err)
		}
	}

	return nil
}
