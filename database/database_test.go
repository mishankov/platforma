package database_test

import (
	"context"
	"testing"

	"github.com/mishankov/platforma/database"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMigrate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("hostamat"),
		postgres.WithUsername("hostamat"),
		postgres.WithPassword("hostamat"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to initialize database: %s", err.Error())
	}

	err = ctr.Snapshot(ctx)
	if err != nil {
		t.Fatalf("failed to create snapshot: %s", err.Error())
	}

	dbURL, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err.Error())
	}

	t.Logf("db connection string: %s", dbURL)

	t.Run("initialize and migrate empty database", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			if err != nil {
				t.Fatalf("failed to restore db: %s", err.Error())
			}
		})

		db, err := database.New(dbURL)
		if err != nil {
			t.Fatalf("failed to initialize database: %s", err.Error())
		}

		if db == nil {
			t.Fatalf("database is nil")
		}

		err = db.Migrate(ctx)
		if err != nil {
			t.Fatalf("failed to migrate database: %s", err.Error())
		}
	})

	t.Run("migrate database with single repository", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			if err != nil {
				t.Fatalf("failed to restore db: %s", err.Error())
			}
		})

		db, err := database.New(dbURL)
		if err != nil {
			t.Fatalf("failed to initialize database: %s", err.Error())
		}

		if db == nil {
			t.Fatalf("database is nil")
		}

		db.RegisterRepository("some_repo", simpleRepo{migrations: []database.Migration{{
			ID:   "init",
			Up:   "CREATE TABLE IF NOT EXISTS simple_repo (id TEXT)",
			Down: "DROP TABLE simple_repo",
		}}})

		err = db.Migrate(ctx)
		if err != nil {
			t.Fatalf("failed to migrate database: %s", err.Error())
		}
	})

	t.Run("migrate database with multiple repositories", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			if err != nil {
				t.Fatalf("failed to restore db: %s", err.Error())
			}
		})

		db, err := database.New(dbURL)
		if err != nil {
			t.Fatalf("failed to initialize database: %s", err.Error())
		}

		if db == nil {
			t.Fatalf("database is nil")
		}

		db.RegisterRepository("some_repo", simpleRepo{migrations: []database.Migration{{
			ID:   "init",
			Up:   "CREATE TABLE IF NOT EXISTS simple_repo (id TEXT)",
			Down: "DROP TABLE simple_repo",
		}}})

		db.RegisterRepository("other_repo", simpleRepo{migrations: []database.Migration{{
			ID:   "init",
			Up:   "CREATE TABLE IF NOT EXISTS other_repo (id TEXT)",
			Down: "DROP TABLE other_repo",
		}}})

		err = db.Migrate(ctx)
		if err != nil {
			t.Fatalf("failed to migrate database: %s", err.Error())
		}
	})

	t.Run("migrate database with failing migration", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			if err != nil {
				t.Fatalf("failed to restore db: %s", err.Error())
			}
		})

		db, err := database.New(dbURL)
		if err != nil {
			t.Fatalf("failed to initialize database: %s", err.Error())
		}

		if db == nil {
			t.Fatalf("database is nil")
		}

		db.RegisterRepository("some_repo", simpleRepo{migrations: []database.Migration{{
			ID:   "init",
			Up:   "CREATE TABLE IF NOT EXISTS simple_repo (id TEXT)",
			Down: "DROP TABLE simple_repo",
		}}})

		db.RegisterRepository("other_repo", simpleRepo{migrations: []database.Migration{{
			ID:   "init",
			Up:   "CREATE TABLE IF NOT EXISTS other_repo (id TEXT)",
			Down: "DROP TABLE other_repo",
		}, {
			ID:   "failing",
			Up:   "not even SQL here",
			Down: "no need for this",
		}}})

		err = db.Migrate(ctx)
		if err == nil {
			t.Fatalf("migration expected to fail")
		}
		t.Logf("migration error: %s", err.Error())
	})
}

type simpleRepo struct {
	migrations []database.Migration
}

func (r simpleRepo) Migrations() []database.Migration {
	return r.migrations
}
