package database_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mishankov/platforma/database"
)

func TestSaveMigrationLogs(t *testing.T) {
	t.Run("single successful", func(t *testing.T) {
		repo := &repoMock{}
		service := database.NewService(repo, &sqlx.DB{})

		err := service.SaveMigrationLogs(context.TODO(), []database.Migration{{ID: "some id"}})
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}
	})
}

type repoMock struct{}

func (r *repoMock) GetMigrationLogs(ctx context.Context) ([]database.MigrationLog, error) {
	return nil, nil
}

func (r *repoMock) SaveMigrationLog(ctx context.Context, log database.MigrationLog) error {
	return nil
}

func (r *repoMock) RemoveMigrationLog(ctx context.Context, repository, id string) error {
	return nil
}

func (r *repoMock) ExecuteQuery(ctx context.Context, query string) error {
	return nil
}

func (r *repoMock) Migrations() []database.Migration {
	return nil
}
