package database

import (
	"time"
)

type MigrationLog struct {
	Repository  string    `db:"repository"`
	MigrationId string    `db:"id"`
	Timestamp   time.Time `db:"timestamp"`
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
