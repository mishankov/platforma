package queue

import (
	"context"
	"sync"
	"time"

	"github.com/platforma-dev/platforma/log"
)

type Handler[T any] interface {
	Handle(ctx context.Context, job T)
}

type HandlerFunc[T any] func(ctx context.Context, job T)

func (f HandlerFunc[T]) Handle(ctx context.Context, job T) {
	f(ctx, job)
}

type Queue[T any] struct {
	jobChan         chan T
	handler         Handler[T]
	wg              sync.WaitGroup
	workersAmount   int
	shutdownTimeout time.Duration
}

func New[T any](handler Handler[T], workersAmount, bufferSize int, shutdownTimeout time.Duration) *Queue[T] {
	return &Queue[T]{jobChan: make(chan T, bufferSize), handler: handler, workersAmount: workersAmount, shutdownTimeout: shutdownTimeout}
}

func (q *Queue[T]) Enqueue(job T) {
	q.jobChan <- job
}

func (q *Queue[T]) Run(ctx context.Context) error {
	q.wg.Add(q.workersAmount)
	for i := range q.workersAmount {
		go q.worker(ctx, i)
	}

	q.wg.Wait()

	return nil
}

func (q *Queue[T]) worker(ctx context.Context, id int) {
	defer q.wg.Done()
	defer log.InfoContext(ctx, "worker finished", "workerId", id)

	log.InfoContext(ctx, "worker started", "workerId", id)

	for {
		breakLoop := false
		select {
		case job := <-q.jobChan:
			q.handler.Handle(ctx, job)
		case <-ctx.Done():
			breakLoop = true
		}

		if breakLoop {
			log.InfoContext(ctx, "shutting down worker")
			break
		}
	}

	timer := time.NewTimer(q.shutdownTimeout)
	defer timer.Stop()

	for {
		select {
		case job := <-q.jobChan:
			q.handler.Handle(ctx, job)
		case <-timer.C:
			return
		}
	}
}
