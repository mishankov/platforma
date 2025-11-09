package queue

import (
	"context"
	"errors"
	"sync"
	"time"
)

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
	select {
	case q.ch <- job:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("timeout")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *ChanQueue[T]) GetJobChan(ctx context.Context) (chan T, error) {
	return q.ch, nil
}
