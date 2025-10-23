package application

import "context"

// Runner is an interface that defines the Run method for executing a task with context.
type Runner interface {
	// Run executes the task with the given context and returns an error if any.
	Run(context.Context) error
}

// RunnerFunc is a function type that implements the Runner interface.
type RunnerFunc func(context.Context) error

// Run executes the RunnerFunc with the given context and returns the result.
func (rf RunnerFunc) Run(ctx context.Context) error {
	return rf(ctx)
}
