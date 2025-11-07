package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/log"

	"github.com/google/uuid"
)

// SchedulerAction defines the function signature for actions executed by the scheduler.
type SchedulerAction func(context.Context) error

// Scheduler represents a periodic task runner that executes an action at fixed intervals.
type Scheduler struct {
	period time.Duration      // The interval between action executions
	runner application.Runner // The runner to execute periodically
}

// New creates a new Scheduler instance with the specified period and action.
func New(period time.Duration, runner application.Runner) *Scheduler {
	return &Scheduler{period: period, runner: runner}
}

// Run starts the scheduler and executes the runner at the configured interval.
// The scheduler will continue running until the context is canceled.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runCtx := context.WithValue(ctx, log.TraceIDKey, uuid.NewString())
			log.InfoContext(runCtx, "scheduler task started")

			err := s.runner.Run(runCtx)
			if err != nil {
				log.ErrorContext(runCtx, "error in scheduler", "error", err)
			}

			log.InfoContext(runCtx, "scheduler task finished")
		case <-ctx.Done():
			return fmt.Errorf("scheduler context canceled: %w", ctx.Err())
		}
	}
}
