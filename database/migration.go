package database

import (
	"database/sql"
	"time"
)

// TODO: rename to migration
type migrations struct {
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

type Schema struct {
	Queries []string
}

type schemer interface {
	Schema() ([]Migration, Schema)
}
