package queue

import (
	"context"
	"fmt"
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

type QueueProvider[T any] interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	EnqueueJob(ctx context.Context, job T) error
	GetJobChan(ctx context.Context) (chan T, error)
}

type QueueProcessor[T any] struct {
	handler         Handler[T]
	queue           QueueProvider[T]
	wg              sync.WaitGroup
	workersAmount   int
	shutdownTimeout time.Duration
}

func New[T any](handler Handler[T], queue QueueProvider[T], workersAmount int, shutdownTimeout time.Duration) *QueueProcessor[T] {
	return &QueueProcessor[T]{handler: handler, queue: queue, workersAmount: workersAmount, shutdownTimeout: shutdownTimeout}
}

func (p *QueueProcessor[T]) Enqueue(ctx context.Context, job T) error {
	err := p.queue.EnqueueJob(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}
	return nil
}

func (p *QueueProcessor[T]) Run(ctx context.Context) error {
	err := p.queue.Open(ctx)
	if err != nil {
		return fmt.Errorf("failed to open queue: %w", err)
	}
	defer func() {
		err := p.queue.Close(ctx)
		if err != nil {
			log.ErrorContext(ctx, "failed to close queue", "error", err)
		}
	}()

	p.wg.Add(p.workersAmount)
	for i := range p.workersAmount {
		go p.worker(ctx, i)
	}

	p.wg.Wait()

	return nil
}

func (p *QueueProcessor[T]) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	defer log.InfoContext(ctx, "worker finished")
	defer func() {
		if r := recover(); r != nil {
			log.ErrorContext(ctx, "worker panic recovered", "panic", r)
		}
	}()

	ctx = context.WithValue(ctx, log.WorkerIDKey, id)

	log.InfoContext(ctx, "worker started")

	jobChan, err := p.queue.GetJobChan(ctx)
	if err != nil {
		log.ErrorContext(ctx, "failed to get job chan", "error", err)
		return
	}

	// we first check for ctx.Done() in separate select statement
	// because select statements choose randomly if both cases are ready
	for {
		breakLoop := false

		select {
		case <-ctx.Done():
			log.InfoContext(ctx, "skipping job due to shutdown")
			breakLoop = true
		default:
			select {
			case job := <-jobChan:
				p.handler.Handle(ctx, job)

			case <-ctx.Done():
				log.InfoContext(ctx, "shutting down worker")
				breakLoop = true
			}
		}

		if breakLoop {
			break
		}
	}

	// after context is cancelled we try to drain remainning jobs from channel
	// before shutdown time expired
	shutdownCtx := context.WithoutCancel(ctx)
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, p.shutdownTimeout)
	defer cancel()

	// same logic with nested select statements as in main loop
	for {
		select {
		case <-shutdownCtx.Done():
			log.InfoContext(shutdownCtx, "shutdown timeout expired")
			return
		default:
			select {
			case job := <-jobChan:
				p.handler.Handle(shutdownCtx, job)
			case <-shutdownCtx.Done():
				log.InfoContext(shutdownCtx, "shutdown timeout expired")
				return
			}
		}
	}
}
