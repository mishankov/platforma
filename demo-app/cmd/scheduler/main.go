package main

import (
	"context"
	"time"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/scheduler"
)

func scheduledTask(ctx context.Context) error {
	log.InfoContext(ctx, "scheduled task executed")
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	s := scheduler.New(time.Second, application.RunnerFunc(scheduledTask))

	go func() {
		time.Sleep(3500 * time.Millisecond)
		cancel()
	}()

	s.Run(ctx)
}
