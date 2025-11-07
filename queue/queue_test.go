package queue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/platforma-dev/platforma/queue"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	type job struct {
		data int
	}

	t.Run("simple queue", func(t *testing.T) {
		t.Parallel()
		res := 0

		q := queue.New(queue.HandlerFunc[job](func(ctx context.Context, job job) {
			res += job.data
		}), 3, 4, time.Microsecond)

		go q.Run(context.TODO())

		q.Enqueue(job{data: 1})
		q.Enqueue(job{data: 1})
		q.Enqueue(job{data: 1})

		time.Sleep(1 * time.Millisecond)

		if res != 3 {
			t.Errorf("expected res to be 1, got %d", res)
		}
	})

	t.Run("shutdown", func(t *testing.T) {
		t.Parallel()
		res := 0
		ctx, cancel := context.WithCancel(context.Background())

		// Defining queue with slow running handler
		q := queue.New(queue.HandlerFunc[job](func(ctx context.Context, job job) {
			time.Sleep(3 * time.Second)
			res += job.data
		}), 1, 10, 1*time.Second)

		// Start queue
		var wg sync.WaitGroup
		wg.Go(func() {
			q.Run(ctx)
		})

		// Enqueue 5 jobs
		for range 5 {
			q.Enqueue(job{data: 1})
		}

		// Wait a bit for workers to pickup jobs
		time.Sleep(1 * time.Second)

		// Cancelling context and waiting for queue to stop
		cancel()
		wg.Wait()

		if res != 2 {
			t.Errorf("expected res to be 2, got %d", res)
		}
	})
}
