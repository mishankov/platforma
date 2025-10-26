package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type service struct {
	repo *repository
	db   *sqlx.DB
}

func newService(repo *repository, db *sqlx.DB) *service {
	return &service{repo: repo, db: db}
}

func (s *service) GetMigrations(ctx context.Context) ([]*migrations, error) {
	return s.repo.GetMigrations(ctx)
}

func (s *service) SaveMigration(ctx context.Context, migration migrations) error {
	return s.repo.SaveMigration(ctx, migration)
}
