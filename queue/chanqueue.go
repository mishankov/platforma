package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrTimeout = errors.New("timeout")
var ErrClosedQueue = errors.New("queue is closed")

type ChanQueue[T any] struct {
	ch         chan T
	mu         sync.Mutex
	opened     bool
	bufferSize int
}

func NewChanQueue[T any](bufferSize int) *ChanQueue[T] {
	return &ChanQueue[T]{bufferSize: bufferSize, opened: false}
}

func (q *ChanQueue[T]) Open(ctx context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.opened {
		q.ch = make(chan T, q.bufferSize)
		q.opened = true
	}

	return nil
}

func (q *ChanQueue[T]) Close(ctx context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.opened {
		close(q.ch)
	}

	return nil
}

func (q *ChanQueue[T]) EnqueueJob(ctx context.Context, job T) error {
	if q.opened {
		select {
		case q.ch <- job:
			return nil
		case <-time.After(5 * time.Second):
			return ErrTimeout
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		}
	}

	return ErrClosedQueue
}

func (q *ChanQueue[T]) GetJobChan(ctx context.Context) (chan T, error) {
	return q.ch, nil
}
