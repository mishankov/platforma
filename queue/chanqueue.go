// Package queue provides a channel-based queue implementation and job processing utilities.
package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrTimeout is returned when an enqueue operation times out.
var ErrTimeout = errors.New("timeout")

// ErrClosedQueue is returned when attempting to operate on a closed queue.
var ErrClosedQueue = errors.New("queue is closed")

// ChanQueue is a thread-safe channel-based queue implementation.
type ChanQueue[T any] struct {
	ch             chan T
	mu             sync.Mutex
	opened         bool
	bufferSize     int
	enqueueTiemout time.Duration
}

// NewChanQueue creates a new channel-based queue with the specified buffer size and enqueue timeout.
func NewChanQueue[T any](bufferSize int, enqueueTimeout time.Duration) *ChanQueue[T] {
	return &ChanQueue[T]{bufferSize: bufferSize, enqueueTiemout: enqueueTimeout, opened: false}
}

// Open initializes the queue and makes it ready to accept jobs.
func (q *ChanQueue[T]) Open(_ context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.opened {
		q.ch = make(chan T, q.bufferSize)
		q.opened = true
	}

	return nil
}

// Close closes the queue and prevents further operations.
func (q *ChanQueue[T]) Close(_ context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.opened {
		close(q.ch)
	}

	return nil
}

// EnqueueJob adds a job to the queue with timeout support.
func (q *ChanQueue[T]) EnqueueJob(ctx context.Context, job T) error {
	if q.opened {
		select {
		case q.ch <- job:
			return nil
		case <-time.After(q.enqueueTiemout):
			return ErrTimeout
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		}
	}

	return ErrClosedQueue
}

// GetJobChan returns the underlying channel for reading jobs.
func (q *ChanQueue[T]) GetJobChan(_ context.Context) (chan T, error) {
	return q.ch, nil
}
