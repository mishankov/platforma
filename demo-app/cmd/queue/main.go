package main

import (
	"context"
	"time"

	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/queue"
)

type job struct {
	data int
}

func jobHandler(ctx context.Context, job job) {
	log.InfoContext(ctx, "job handled", "data", job.data)
}

func main() {
	ctx := context.Background()

	q := queue.NewChanQueue[job](10, 3*time.Second)
	p := queue.New(queue.HandlerFunc[job](jobHandler), q, 2, time.Second)

	go p.Run(ctx)
	time.Sleep(time.Millisecond)

	p.Enqueue(ctx, job{data: 1})
	p.Enqueue(ctx, job{data: 2})
	p.Enqueue(ctx, job{data: 3})

	time.Sleep(time.Millisecond)
}
