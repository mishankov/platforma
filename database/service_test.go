package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mishankov/platforma/database"
)

func TestSaveMigrationLogs(t *testing.T) {
	t.Parallel()
	t.Run("single successful", func(t *testing.T) {
		t.Parallel()
		repo := &repoMock{}
		service := database.NewService(repo, &sqlx.DB{})

		err := service.SaveMigrationLogs(context.TODO(), []database.Migration{{ID: "some id"}})
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}
	})

	t.Run("single error", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		repo := &repoMock{saveMigrationLog: func(ctx context.Context, ml database.MigrationLog) error {
			return ErrSome
		}}
		service := database.NewService(repo, &sqlx.DB{})

		err := service.SaveMigrationLogs(context.TODO(), []database.Migration{{ID: "some id"}})
		if err == nil {
			t.Fatalf("expected error, got nothing")
		}

		if !errors.Is(err, ErrSome) {
			t.Fatalf("expected ErrSome, got: %s", err.Error())
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		ErrOther := errors.New("other error")
		repo := &repoMock{saveMigrationLog: func(ctx context.Context, ml database.MigrationLog) error {
			if ml.MigrationId == "some id" {
				return ErrSome
			}

			if ml.MigrationId == "other id" {
				return ErrOther
			}

			return nil
		}}
		service := database.NewService(repo, &sqlx.DB{})

		err := service.SaveMigrationLogs(context.TODO(), []database.Migration{{ID: "some id"}, {ID: "other id"}})
		if err == nil {
			t.Fatalf("expected error, got nothing")
		}

		if !errors.Is(err, ErrSome) {
			t.Fatalf("expected ErrSome, got: %s", err.Error())
		}

		if !errors.Is(err, ErrOther) {
			t.Fatalf("expected ErrOther, got: %s", err.Error())
		}
	})
}

type repoMock struct {
	saveMigrationLog func(context.Context, database.MigrationLog) error
}

func (r *repoMock) GetMigrationLogs(ctx context.Context) ([]database.MigrationLog, error) {
	return nil, nil
}

func (r *repoMock) SaveMigrationLog(ctx context.Context, log database.MigrationLog) error {
	if r.saveMigrationLog != nil {
		return r.saveMigrationLog(ctx, log)
	}

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
