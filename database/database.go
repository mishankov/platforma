package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
	repositories map[string]any
	migrators    map[string]migrator
	repository   *Repository
	service      *service
}

func New(connection string) (*Database, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repository := newRepository(db)
	service := newService(repository)
	return &Database{DB: db, repositories: make(map[string]any), migrators: make(map[string]migrator), repository: repository, service: service}, nil
}

func (db *Database) RegisterRepository(name string, repository any) {
	db.repositories[name] = repository

	if migr, ok := repository.(migrator); ok {
		db.migrators[name] = migr
	}
}

func (db *Database) Migrate(ctx context.Context) error {
	// Ensure that migration table exists
	err := db.service.MigrateSelf(ctx)
	if err != nil {
		return err
	}

	// Get completed migrations
	migrationLogs, err := db.service.GetMigrationLogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to select migrations state: %w", err)
	}

	// Get migrations from all migrators
	migrations := []Migration{}
	for name, migrator := range db.migrators {
		for _, migr := range migrator.Migrations() {
			migr.repository = name
			migrations = append(migrations, migr)
		}
	}

	err = db.service.ApplyMigrations(ctx, migrations, migrationLogs)
	if err != nil {
		return err
	}

	return nil
}
