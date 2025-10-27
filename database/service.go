package database

import (
	"context"
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/mishankov/platforma/log"
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

func (s *service) MigrateSelf(ctx context.Context) error {
	migrations := s.repo.Migrations()
	appliedMigrations := []Migration{}
	migrationLogs, err := s.repo.GetMigrationLogs(ctx)

	// If GetMigrationLogs returns error, log table probably does not exist,
	// so we should apply all migrations for it
	if err != nil {
		for _, migr := range migrations {
			err := s.ApplyMigration(ctx, migr)
			if err != nil {
				s.RevertMigrations(ctx, appliedMigrations)
				return err
			}
			appliedMigrations = append(appliedMigrations, migr)
		}
	}

	for _, migr := range migrations {
		if !slices.ContainsFunc(migrationLogs, func(l migrationLog) bool {
			return l.Repository == "platforma_migrations" && l.MigrationId.String == migr.ID
		}) {
			err := s.ApplyMigration(ctx, migr)
			if err != nil {
				s.RevertMigrations(ctx, appliedMigrations)
				return err
			}
		}
	}

	return nil
}

func (s *service) ApplyMigration(ctx context.Context, migration Migration) error {
	_, err := s.db.ExecContext(ctx, migration.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}
	return nil
}

func (s *service) ApplyMigrations(ctx context.Context, migrations []Migration, migrationLogs []migrationLog) error {
	appliedMigrations := []Migration{}
	for _, migr := range migrations {
		if !slices.ContainsFunc(migrationLogs, func(l migrationLog) bool {
			return l.Repository == migr.repository && l.MigrationId.String == migr.ID
		}) {
			err := s.ApplyMigration(ctx, migr)
			if err != nil {
				s.RevertMigrations(ctx, appliedMigrations)
				return err
			}
		}
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

func (s *service) RevertMigrations(ctx context.Context, migrations []Migration) {
	for _, migr := range slices.Backward(migrations) {
		err := s.RevertMigration(ctx, migr)
		if err != nil {
			log.ErrorContext(ctx, "failed to revert migration", "migration", migr.ID, "error", err.Error())
		}
	}
}
