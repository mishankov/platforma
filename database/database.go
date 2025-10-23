package database

import (
	"context"
	"slices"
	"time"

	"github.com/mishankov/platforma/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
	repositories map[string]any
	migrators    map[string]shemer
}

func New(connection string) (*Database, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db, repositories: make(map[string]any), migrators: make(map[string]shemer)}, nil
}

func (db *Database) RegisterRepository(name string, repository any) {
	db.repositories[name] = repository

	if migr, ok := repository.(shemer); ok {
		db.migrators[name] = migr
	}
}

func (db *Database) Migrate(ctx context.Context) error {
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS platforma_migrations (repository TEXT, id TEXT, timestamp TIMESTAMP)"); err != nil {
		return err
	}

	// Select data from platforma_migrations table
	var migrationsState []migrations
	err := db.SelectContext(ctx, &migrationsState, "SELECT * FROM platforma_migrations")
	if err != nil {
		return err
	}

	appliedMigrations := []Migration{}
	migrationErr := error(nil)

	for repoName, migr := range db.migrators {
		repoMigrations, repoSchema := migr.Schema()

		repoHasMigrations := slices.ContainsFunc(migrationsState, func(m migrations) bool {
			return m.Repository == repoName
		})

		// If repo does not has migrations apply schema and exit
		if !repoHasMigrations {
			for _, query := range repoSchema.Queries {
				if _, err := db.ExecContext(ctx, query); err != nil {
					migrationErr = err
					break
				}
				log.InfoContext(ctx, "schema applied", "repository", repoName)
			}

			// Log that schema applied
			if _, err := db.ExecContext(ctx, "INSERT INTO platforma_migrations (repository, timestamp) VALUES ($1, $2)", repoName, time.Now()); err != nil {
				return err
			}

			// If schema is applied, log that all migrations are also applied
			for _, migration := range repoMigrations {
				if _, err := db.ExecContext(ctx, "INSERT INTO platforma_migrations (repository, id, timestamp) VALUES ($1, $2, $3)", repoName, migration.ID, time.Now()); err != nil {
					return err
				}
			}

			continue
		}

		for _, migration := range repoMigrations {
			migration.repository = repoName

			// Check if migration has been applied
			migrationHasApplied := slices.ContainsFunc(migrationsState, func(m migrations) bool {
				return m.Repository == repoName && m.MigrationId.String == migration.ID
			})

			if migrationHasApplied {
				continue
			}

			if _, err := db.ExecContext(ctx, migration.Up); err != nil {
				migrationErr = err
				log.ErrorContext(ctx, "failed to apply migration for repository", "migration", migration.ID, "repository", repoName)
				break
			}

			appliedMigrations = append(appliedMigrations, migration)
			log.InfoContext(ctx, "applied migration for repository", "migration", migration.ID, "repository", repoName)

			// Log that migration applied
			if _, err := db.ExecContext(ctx, "INSERT INTO platforma_migrations (repository, id, timestamp) VALUES ($1, $2, $3)", repoName, migration.ID, time.Now()); err != nil {
				return err
			}
		}

		if migrationErr != nil {
			break
		}
	}

	if migrationErr != nil {
		for _, migration := range slices.Backward(appliedMigrations) {
			if _, err := db.ExecContext(ctx, migration.Down); err != nil {
				log.ErrorContext(ctx, "failed to rollback migration %s for repository %s", migration.ID, migration.repository)
				return err
			}

			if _, err := db.ExecContext(ctx, "DELETE FROM platforma_migrations WHERE repository = $1 AND id = $2", migration.repository, migration.ID); err != nil {
				return err
			}
		}
	}

	return migrationErr
}
