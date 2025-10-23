package application

import "fmt"

// ErrStartupTaskFailed represents an error that occurs when a startup task fails.
type ErrStartupTaskFailed struct {
	err error
}

// Error returns the formatted error message for ErrStartupTaskFailed.
func (e *ErrStartupTaskFailed) Error() string {
	return fmt.Sprintf("failed to run startup task: %v", e.err)
}

// Unwrap returns the underlying error for ErrStartupTaskFailed.
func (e *ErrStartupTaskFailed) Unwrap() error {
	return e.err
}

// StartupTaskConfig contains configuration options for a startup task.
type StartupTaskConfig struct {
	Name         string // Name of the startup task
	AbortOnError bool   // Whether to abort application startup if this task fails
}

// startupTask represents an individual startup task with its runner and configuration.
type startupTask struct {
	runner Runner
	config StartupTaskConfig
}
