package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/platforma-dev/platforma/auth"
)

func TestDeleteUser_WithCleanupEnqueuer_EnqueuesJob(t *testing.T) {
	t.Parallel()

	mockRepo := &mockRepository{}
	mockAuthStorage := &mockAuthStorage{}
	mockEnqueuer := &mockCleanupEnqueuer{}

	service := auth.NewService(mockRepo, mockAuthStorage, "session", nil, nil, mockEnqueuer)

	user := &auth.User{ID: "test-user-id"}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, user)

	err := service.DeleteUser(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mockEnqueuer.enqueueCalled {
		t.Fatal("expected Enqueue to be called")
	}

	if mockEnqueuer.lastJob.UserID != "test-user-id" {
		t.Fatalf("expected UserID 'test-user-id', got %q", mockEnqueuer.lastJob.UserID)
	}

	if mockEnqueuer.lastJob.DeletedAt.IsZero() {
		t.Fatal("expected DeletedAt to be set")
	}
}

func TestDeleteUser_WithNilEnqueuer_Succeeds(t *testing.T) {
	t.Parallel()

	mockRepo := &mockRepository{}
	mockAuthStorage := &mockAuthStorage{}

	service := auth.NewService(mockRepo, mockAuthStorage, "session", nil, nil, nil)

	user := &auth.User{ID: "test-user-id"}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, user)

	err := service.DeleteUser(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteUser_EnqueueError_StillSucceeds(t *testing.T) {
	t.Parallel()

	mockRepo := &mockRepository{}
	mockAuthStorage := &mockAuthStorage{}
	mockEnqueuer := &mockCleanupEnqueuer{
		enqueueErr: errors.New("queue is full"),
	}

	service := auth.NewService(mockRepo, mockAuthStorage, "session", nil, nil, mockEnqueuer)

	user := &auth.User{ID: "test-user-id"}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, user)

	err := service.DeleteUser(ctx)
	if err != nil {
		t.Fatalf("expected no error even with enqueue failure, got %v", err)
	}

	if !mockEnqueuer.enqueueCalled {
		t.Fatal("expected Enqueue to be called")
	}
}

type mockCleanupEnqueuer struct {
	enqueueCalled bool
	lastJob       auth.UserCleanupJob
	enqueueErr    error
}

func (m *mockCleanupEnqueuer) Enqueue(_ context.Context, job auth.UserCleanupJob) error {
	m.enqueueCalled = true
	m.lastJob = job
	return m.enqueueErr
}

type mockRepository struct{}

func (m *mockRepository) Get(_ context.Context, _ string) (*auth.User, error) {
	return nil, nil
}

func (m *mockRepository) GetByUsername(_ context.Context, _ string) (*auth.User, error) {
	return nil, nil
}

func (m *mockRepository) Create(_ context.Context, _ *auth.User) error {
	return nil
}

func (m *mockRepository) UpdatePassword(_ context.Context, _, _, _ string) error {
	return nil
}

func (m *mockRepository) Delete(_ context.Context, _ string) error {
	return nil
}

type mockAuthStorage struct{}

func (m *mockAuthStorage) GetUserIdFromSessionId(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (m *mockAuthStorage) CreateSessionForUser(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (m *mockAuthStorage) DeleteSession(_ context.Context, _ string) error {
	return nil
}

func (m *mockAuthStorage) DeleteSessionsByUserId(_ context.Context, _ string) error {
	return nil
}
