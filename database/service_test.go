package database

import (
	"context"
	"errors"
	"testing"
)

func TestSaveMigrationLogs(t *testing.T) {
	t.Parallel()
	t.Run("single successful", func(t *testing.T) {
		t.Parallel()
		repo := &repoMock{}
		service := newService(repo)

		err := service.SaveMigrationLogs(context.TODO(), []Migration{{ID: "some id"}})
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}
	})

	t.Run("single error", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		repo := &repoMock{saveMigrationLog: func(ctx context.Context, ml MigrationLog) error {
			return ErrSome
		}}
		service := newService(repo)

		err := service.SaveMigrationLogs(context.TODO(), []Migration{{ID: "some id"}})
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
		repo := &repoMock{saveMigrationLog: func(ctx context.Context, ml MigrationLog) error {
			if ml.MigrationId == "some id" {
				return ErrSome
			}

			if ml.MigrationId == "other id" {
				return ErrOther
			}

			return nil
		}}
		service := newService(repo)

		err := service.SaveMigrationLogs(context.TODO(), []Migration{{ID: "some id"}, {ID: "other id"}})
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
		service := newService(repo)

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
		repo := &repoMock{getMigrationLogs: func(ctx context.Context) ([]MigrationLog, error) {
			return []MigrationLog{{}, {}}, nil
		}}
		service := newService(repo)

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
		repo := &repoMock{getMigrationLogs: func(ctx context.Context) ([]MigrationLog, error) {
			return nil, ErrSome
		}}
		service := newService(repo)

		_, err := service.GetMigrationLogs(context.TODO())
		if err == nil {
			t.Fatalf("expected error, got nothing")
		}

		if !errors.Is(err, ErrSome) {
			t.Fatalf("expected ErrSome, got: %s", err.Error())
		}
	})
}

func TestRevertMigrations(t *testing.T) {
	t.Parallel()
	t.Run("successful revert", func(t *testing.T) {
		t.Parallel()
		repo := &repoMock{}
		service := newService(repo)

		migrations := []Migration{{ID: "migration1"}, {ID: "migration2"}}

		err := service.RevertMigrations(context.TODO(), migrations)
		if err != nil {
			t.Fatalf("expected no errors, got: %s", err.Error())
		}
	})

	t.Run("revert with error", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		repo := &repoMock{
			executeQuery: func(ctx context.Context, query string) error {
				return ErrSome
			},
		}
		service := newService(repo)

		err := service.RevertMigrations(context.TODO(), []Migration{{ID: "migration1"}})
		if err == nil {
			t.Fatalf("expected error, got nothing")
		}

		if !errors.Is(err, ErrSome) {
			t.Fatalf("expected ErrSome, got: %s", err.Error())
		}
	})

	t.Run("revert with multiple errors", func(t *testing.T) {
		t.Parallel()
		ErrSome := errors.New("some error")
		ErrOther := errors.New("other error")
		repo := &repoMock{
			executeQuery: func(ctx context.Context, query string) error {
				if query == "migration1" {
					return ErrSome
				}

				if query == "migration2" {
					return ErrOther
				}

				return nil
			},
		}
		service := newService(repo)

		err := service.RevertMigrations(context.TODO(), []Migration{{Down: "migration1"}, {Down: "migration2"}})
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
	getMigrationLogs func(context.Context) ([]MigrationLog, error)
	saveMigrationLog func(context.Context, MigrationLog) error
	executeQuery     func(context.Context, string) error
}

func (r *repoMock) GetMigrationLogs(ctx context.Context) ([]MigrationLog, error) {
	if r.getMigrationLogs != nil {
		return r.getMigrationLogs(ctx)
	}
	return nil, nil
}

func (r *repoMock) SaveMigrationLog(ctx context.Context, log MigrationLog) error {
	if r.saveMigrationLog != nil {
		return r.saveMigrationLog(ctx, log)
	}

	return nil
}

func (r *repoMock) RemoveMigrationLog(ctx context.Context, repository, id string) error {
	return nil
}

func (r *repoMock) ExecuteQuery(ctx context.Context, query string) error {
	if r.executeQuery != nil {
		return r.executeQuery(ctx, query)
	}
	return nil
}

func (r *repoMock) Migrations() []Migration {
	return nil
}
