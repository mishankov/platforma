package auth

import (
	"context"
	"time"
)

// CleanupEnqueuer defines the interface for enqueuing user cleanup jobs.
// Applications can implement this interface to handle cleanup tasks after user deletion,
// such as removing user files, canceling subscriptions, or anonymizing data.
//
// This interface is designed to work with queue.Processor[UserCleanupJob] from the
// platforma queue package, allowing cleanup jobs to be processed asynchronously.
type CleanupEnqueuer interface {
	Enqueue(ctx context.Context, job UserCleanupJob) error
}

// UserCleanupJob represents a cleanup job that is enqueued after a user is deleted.
// It contains the necessary information for cleanup handlers to identify and process
// the cleanup tasks associated with the deleted user.
type UserCleanupJob struct {
	// UserID is the unique identifier of the deleted user.
	UserID string `json:"userId"`
	// DeletedAt is the timestamp when the user was deleted.
	DeletedAt time.Time `json:"deletedAt"`
}
