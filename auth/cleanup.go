package auth

import (
	"context"
	"time"
)

type cleanupEnqueuer interface {
	Enqueue(ctx context.Context, job UserCleanupJob) error
}

// UserCleanupJob represents a cleanup job that is enqueued after a user is deleted.
// It contains the necessary information for cleanup handlers to identify and process
// the cleanup tasks associated with the deleted user.
type UserCleanupJob struct {
	UserID    string
	DeletedAt time.Time
}
