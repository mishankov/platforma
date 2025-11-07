package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/platforma-dev/platforma/queue"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	type job struct {
		data int
	}
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
}
