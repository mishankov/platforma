package database

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/mishankov/platforma/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
	repositories map[string]any
	migrators    map[string]schemer
	repository   *repository
	service      *service
}

func New(connection string) (*Database, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repository := newRepository(db)
	service := newService(repository, db)
	return &Database{DB: db, repositories: make(map[string]any), migrators: make(map[string]schemer), repository: repository, service: service}, nil
}

func (db *Database) RegisterRepository(name string, repository any) {
	db.repositories[name] = repository

	if migr, ok := repository.(schemer); ok {
		db.migrators[name] = migr
	}
}

func (db *Database) Migrate(ctx context.Context) error {
	// Ensure that migration table exists
	_, databaseRepoSchema := db.repository.Schema()
	err := db.service.ApplySchema(ctx, databaseRepoSchema)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get completed migrations
	migrationsState, err := db.service.GetMigrationLogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to select migrations state: %w", err)
	}

	appliedMigrations := []Migration{}
	migrationErr := error(nil)

	for repoName, migr := range db.migrators {
		repoMigrations, repoSchema := migr.Schema()

		repoHasMigrations := slices.ContainsFunc(migrationsState, func(m migrationLog) bool {
			return m.Repository == repoName
		})

		// If repo does not has migrations apply schema and exit
		if !repoHasMigrations {
			err := db.service.ApplySchema(ctx, repoSchema)
			if err != nil {
				return fmt.Errorf("failed to execute schema query: %w", err)
			}
			log.InfoContext(ctx, "schema applied", "repository", repoName)

			// Log that schema applied
			err = db.service.SaveMigrationLog(ctx, migrationLog{Repository: repoName, Timestamp: time.Now()})
			if err != nil {
				return fmt.Errorf("failed to insert migration record: %w", err)
			}

			// If schema is applied, log that all migrations are also applied
			for _, migration := range repoMigrations {
				err := db.service.SaveMigrationLog(ctx, migrationLog{Repository: repoName, MigrationId: sql.NullString{String: migration.ID}, Timestamp: time.Now()})
				if err != nil {
					return fmt.Errorf("failed to insert migration record: %w", err)
				}
			}

			continue
		}

		for _, migration := range repoMigrations {
			migration.repository = repoName

			// Check if migration has been applied
			migrationHasApplied := slices.ContainsFunc(migrationsState, func(m migrationLog) bool {
				return m.Repository == repoName && m.MigrationId.String == migration.ID
			})

			// Skip to next migration if current is applied
			if migrationHasApplied {
				continue
			}

			// Try to apply mifration
			err := db.service.ApplyMigration(ctx, migration)
			if err != nil {
				migrationErr = fmt.Errorf("failed to apply migration %s for repository %s: %w", migration.ID, repoName, err)
				break
			}

			// If migration applied successfuly add it to applied migrations list
			appliedMigrations = append(appliedMigrations, migration)
			log.InfoContext(ctx, "applied migration for repository", "migration", migration.ID, "repository", repoName)

			// Log that migration applied
			err = db.service.SaveMigrationLog(ctx, migrationLog{Repository: repoName, MigrationId: sql.NullString{String: migration.ID}, Timestamp: time.Now()})
			if err != nil {
				migrationErr = fmt.Errorf("failed to insert migration record: %w", err)
				break
			}
		}

		if migrationErr != nil {
			break
		}
	}

	if migrationErr != nil {
		for _, migration := range slices.Backward(appliedMigrations) {
			err := db.service.RevertMigration(ctx, migration)
			if err != nil {
				return fmt.Errorf("failed to rollback migration %s for repository %s: %w", migration.ID, migration.repository, err)
			}

			err = db.service.RemoveMigrationLog(ctx, migration.repository, migration.ID)
			if err != nil {
				return fmt.Errorf("failed to delete migration record: %w", err)
			}
		}
	}

	return migrationErr
}
