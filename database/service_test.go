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

func TestGetMigrationLogs(t *testing.T) {
	t.Parallel()
	t.Run("successful no logs", func(t *testing.T) {
		t.Parallel()
		repo := &repoMock{}
		service := database.NewService(repo, &sqlx.DB{})

		logs, err := service.GetMigrationLogs(context.TODO())
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}

		if len(logs) > 0 {
			t.Fatalf("expected no logs, got: %d", len(logs))
		}
	})

	t.Run("successful some logs", func(t *testing.T) {
		t.Parallel()
		repo := &repoMock{getMigrationLogs: func(ctx context.Context) ([]database.MigrationLog, error) {
			return []database.MigrationLog{{}, {}}, nil
		}}
		service := database.NewService(repo, &sqlx.DB{})

		logs, err := service.GetMigrationLogs(context.TODO())
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}

		if len(logs) != 2 {
			t.Fatalf("expected 2 logs, got: %d", len(logs))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		repo := &repoMock{getMigrationLogs: func(ctx context.Context) ([]database.MigrationLog, error) {
			return nil, ErrSome
		}}
		service := database.NewService(repo, &sqlx.DB{})

		_, err := service.GetMigrationLogs(context.TODO())
		if err == nil {
			t.Fatalf("expected error, got nothing")
		}

		if !errors.Is(err, ErrSome) {
			t.Fatalf("expected ErrSome, got: %s", err.Error())
		}
	})
}

type repoMock struct {
	getMigrationLogs func(context.Context) ([]database.MigrationLog, error)
	saveMigrationLog func(context.Context, database.MigrationLog) error
}

func (r *repoMock) GetMigrationLogs(ctx context.Context) ([]database.MigrationLog, error) {
	if r.getMigrationLogs != nil {
		return r.getMigrationLogs(ctx)
	}
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
