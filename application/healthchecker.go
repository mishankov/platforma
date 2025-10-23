package application

import "context"

// Healthchecker represents a service that can perform health checks.
type Healthchecker interface {
	// Healthcheck returns the health status of the service.
	Healthcheck(context.Context) any
}
