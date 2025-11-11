package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/platforma-dev/platforma/queue"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	t.Run("simple queue", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		res := 0

		q := &mockQueue[job]{
			jobChan: make(chan job, 10),
		}

		p := queue.New(queue.HandlerFunc[job](func(ctx context.Context, job job) {
			res += job.data
		}), q, 4, time.Microsecond)

		go p.Run(context.TODO())

		p.Enqueue(ctx, job{data: 1})
		p.Enqueue(ctx, job{data: 1})
		p.Enqueue(ctx, job{data: 1})

		time.Sleep(1 * time.Millisecond)

		if res != 3 {
			t.Errorf("expected res to be 1, got %d", res)
		}
	})
}

type job struct {
	data int
}

type mockQueue[T any] struct {
	jobChan chan T
}

func (q *mockQueue[T]) Open(ctx context.Context) error {
	return nil
}

func (q *mockQueue[T]) Close(ctx context.Context) error {
	return nil
}

func (q *mockQueue[T]) EnqueueJob(ctx context.Context, job T) error {
	q.jobChan <- job
	return nil
}

func (q *mockQueue[T]) GetJobChan(ctx context.Context) (chan T, error) {
	return q.jobChan, nil
}
