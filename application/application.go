package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mishankov/platforma/database"
	"github.com/mishankov/platforma/log"
)

// ErrDatabaseMigrationFailed is an error type that represents a failed database migration.
type ErrDatabaseMigrationFailed struct {
	err error
}

// Error returns the formatted error message for ErrDatabaseMigrationFailed.
func (e *ErrDatabaseMigrationFailed) Error() string {
	return fmt.Sprintf("failed to migrate database: %v", e.err)
}

// Unwrap returns the underlying error for ErrDatabaseMigrationFailed.
func (e *ErrDatabaseMigrationFailed) Unwrap() error {
	return e.err
}

// Application manages startup tasks and services for the application lifecycle.
type Application struct {
	startupTasks   []startupTask
	services       map[string]Runner
	healthcheckers map[string]Healthchecker
	databases      map[string]*database.Database
	health         *ApplicationHealth
}

// New creates and returns a new Application instance.
func New() *Application {
	return &Application{services: make(map[string]Runner), healthcheckers: make(map[string]Healthchecker), databases: make(map[string]*database.Database), health: NewApplicationHealth()}
}

// Health returns the current health status of the application.
func (a *Application) Health(ctx context.Context) *ApplicationHealth {
	for hcName, hc := range a.healthcheckers {
		a.health.SetServiceData(hcName, hc.Healthcheck(ctx))
	}
	return a.health
}

// OnStart registers a new startup task with the given runner and configuration.
func (a *Application) OnStart(task Runner, config StartupTaskConfig) {
	a.startupTasks = append(a.startupTasks, startupTask{task, config})
}

func (a *Application) OnStartFunc(task RunnerFunc, config StartupTaskConfig) {
	a.startupTasks = append(a.startupTasks, startupTask{task, config})
}

// RegisterDatabase adds a database to the application.
func (a *Application) RegisterDatabase(dbName string, db *database.Database) {
	a.databases[dbName] = db
}

// RegisterRepository adds a repository to the application.
func (a *Application) RegisterRepository(dbName string, repoName string, repository any) {
	a.databases[dbName].RegisterRepository(repoName, repository)
}

// RegisterService adds a named service to the application.
func (a *Application) RegisterService(serviceName string, service Runner) {

	a.services[serviceName] = service
	a.health.Services[serviceName] = &ServiceHealth{Status: ServiceStatusNotStarted}

	healthcheckerService, ok := service.(Healthchecker)
	if ok {
		a.healthcheckers[serviceName] = healthcheckerService
	}
}

func (a *Application) RegisterDomain(name, dbName string, domain Domain) {
	if dbName != "" {
		repository := domain.GetRepository()
		a.RegisterRepository(dbName, name+"_repository", repository)
	}
}

// Run executes all startup tasks and services in the application.
// It returns an error if any startup task configured to abort on error fails.
func (a *Application) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	log.InfoContext(ctx, "starting application", "startupTasks", len(a.startupTasks))

	for dbName, db := range a.databases {
		log.InfoContext(ctx, "migrating database", "database", dbName)
		err := db.Migrate(ctx)
		if err != nil {
			log.ErrorContext(ctx, "error in database migration", "error", err, "database", dbName)
			return &ErrDatabaseMigrationFailed{err: err}
		}
	}

	for i, task := range a.startupTasks {
		log.InfoContext(ctx, "running task", "task", task.config.Name, "index", i)

		taskCtx := context.WithValue(ctx, log.StartupTaskKey, task.config.Name)

		err := task.runner.Run(taskCtx)
		if err != nil {
			log.ErrorContext(ctx, "error in startup task", "error", err, "task", task.config.Name)

			if task.config.AbortOnError {
				return &ErrStartupTaskFailed{err: err}
			}
		}
	}

	var wg sync.WaitGroup

	for serviceName, service := range a.services {
		wg.Add(1)

		serviceCtx := context.WithValue(ctx, log.ServiceNameKey, serviceName)

		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.ErrorContext(serviceCtx, "service panicked", "service", serviceName, "panic", r)
				}
			}()

			log.InfoContext(ctx, "starting service", "service", serviceName)
			a.health.StartService(serviceName)

			err := service.Run(serviceCtx)
			if err != nil {
				a.health.FailService(serviceName, err)
				log.ErrorContext(ctx, "error in service", "service", serviceName, "error", err)
			}
		}()
	}

	a.health.StartedAt = time.Now()

	wg.Wait()

	return nil
}
