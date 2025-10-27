package database

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mishankov/platforma/log"
)

type repository interface {
	GetMigrationLogs(ctx context.Context) ([]MigrationLog, error)
	SaveMigrationLog(ctx context.Context, log MigrationLog) error
	RemoveMigrationLog(ctx context.Context, repository, id string) error
	ExecuteQuery(ctx context.Context, query string) error
	Migrations() []Migration
}

type service struct {
	repo repository
}

func NewService(repo repository, db *sqlx.DB) *service {
	return &service{repo: repo}
}

func (s *service) GetMigrationLogs(ctx context.Context) ([]MigrationLog, error) {
	logs, err := s.repo.GetMigrationLogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration logs: %w", err)
	}
	return logs, nil
}

func (s *service) SaveMigrationLog(ctx context.Context, repository, migrationId string) error {
	err := s.repo.SaveMigrationLog(ctx, MigrationLog{Repository: repository, MigrationId: migrationId, Timestamp: time.Now()})
	if err != nil {
		return fmt.Errorf("failed to save migration log: %w", err)
	}
	return nil
}

func (s *service) SaveMigrationLogs(ctx context.Context, migrations []Migration) error {
	masterErr := error(nil)
	for _, migr := range migrations {
		err := s.SaveMigrationLog(ctx, migr.repository, migr.ID)
		if err != nil {
			masterErr = errors.Join(masterErr, err)
		}
	}

	return masterErr
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
				revertErr := s.RevertMigrations(ctx, appliedMigrations)
				if revertErr != nil {
					log.ErrorContext(ctx, "got error(s) trying to revert migrations: %s", revertErr.Error())
				}
				return err
			}
			appliedMigrations = append(appliedMigrations, migr)
		}
	}

	for _, migr := range migrations {
		if !slices.ContainsFunc(migrationLogs, func(l MigrationLog) bool {
			return l.Repository == "platforma_migrations" && l.MigrationId == migr.ID
		}) {
			err := s.ApplyMigration(ctx, migr)
			if err != nil {
				revertErr := s.RevertMigrations(ctx, appliedMigrations)
				if revertErr != nil {
					log.ErrorContext(ctx, "got error(s) trying to revert migrations: %s", revertErr.Error())
				}
				return err
			}
			appliedMigrations = append(appliedMigrations, migr)
		}
	}

	err = s.SaveMigrationLogs(ctx, appliedMigrations)
	if err != nil {
		log.ErrorContext(ctx, "got error(s) trying to save migration logs", "error", err.Error())
	}

	return nil
}

func (s *service) ApplyMigration(ctx context.Context, migration Migration) error {
	err := s.repo.ExecuteQuery(ctx, migration.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}
	return nil
}

func (s *service) ApplyMigrations(ctx context.Context, migrations []Migration, migrationLogs []MigrationLog) error {
	appliedMigrations := []Migration{}
	for _, migr := range migrations {
		if !slices.ContainsFunc(migrationLogs, func(l MigrationLog) bool {
			return l.Repository == migr.repository && l.MigrationId == migr.ID
		}) {
			err := s.ApplyMigration(ctx, migr)
			if err != nil {
				revertErr := s.RevertMigrations(ctx, appliedMigrations)
				if revertErr != nil {
					log.ErrorContext(ctx, "got error(s) trying to revert migrations: %s", revertErr.Error())
				}
				return err
			}
		}
	}

	err := s.SaveMigrationLogs(ctx, appliedMigrations)
	if err != nil {
		log.ErrorContext(ctx, "got error(s) trying to save migration logs", "error", err.Error())
	}

	return nil
}

func (s *service) RevertMigration(ctx context.Context, migration Migration) error {
	err := s.repo.ExecuteQuery(ctx, migration.Down)
	if err != nil {
		return fmt.Errorf("failed to revert migration: %w", err)
	}
	return nil
}

func (s *service) RevertMigrations(ctx context.Context, migrations []Migration) error {
	masterErr := error(nil)
	for _, migr := range slices.Backward(migrations) {
		err := s.RevertMigration(ctx, migr)
		if err != nil {
			masterErr = errors.Join(masterErr, fmt.Errorf("failed to revert migration %s: %w", migr.ID, err))
		}
	}

	return masterErr
}
