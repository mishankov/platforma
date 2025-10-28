package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mishankov/platforma/application"
	"github.com/mishankov/platforma/log"
)

type Clock struct{}

func (r *Clock) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.InfoContext(ctx, "tick")
		case <-ctx.Done():
			log.InfoContext(ctx, "finished")
			return fmt.Errorf("context error: %w", ctx.Err())
		}
	}
}

func main() {
	ctx := context.Background()
	app := application.New()

	app.RegisterService("clock", &Clock{})

	if err := app.Run(ctx); err != nil {
		log.ErrorContext(ctx, "app finished with error", "error", err)
	}
}
