package database

import (
	"database/sql"
	"time"
)

type migrationLog struct {
	Repository  string         `db:"repository"`
	MigrationId sql.NullString `db:"id"`
	Timestamp   time.Time      `db:"timestamp"`
}

type Migration struct {
	ID         string
	Up         string
	Down       string
	repository string
}

type migrator interface {
	Migrations() []Migration
}
