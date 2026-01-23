# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Platforma is a Go framework for building web applications with built-in support for HTTP servers, database migrations, authentication, background job processing, and scheduled tasks. It uses a layered architecture with Application as the central orchestrator.

## Build and Test Commands

### Using Task (Taskfile)
- **Run linter**: `task lint`
- **Auto-fix linting issues**: `task fix`
- **Run tests with coverage**: `task test` (generates `coverage.html` and `coverage.out`)
- **Run both lint and test**: `task check`
- **Run documentation server**: `task docs` (runs in `docs/` directory)

### Using Go directly
- **Run all tests**: `go test ./...`
- **Run tests with coverage**: `go test -coverprofile=coverage.tmp.out ./...`
- **Build demo apps**: `go build ./demo-app/cmd/...`

### CLI Commands
The main binary provides CLI commands:
- **Version**: `go run main.go --version` or `go run main.go -v`
- **Generate**: `go run main.go generate [args]`
- **Docs**: `go run main.go docs [args]`

## Architecture

### Core Application Pattern

The framework follows a lifecycle-based architecture centered around `application.Application`:

1. **Application** (`application/application.go`) - Central orchestrator that manages:
   - **Startup Tasks**: One-time initialization tasks that run sequentially before services start
   - **Services**: Long-running components (HTTP servers, queue processors, schedulers) that run concurrently
   - **Databases**: PostgreSQL databases with automatic migration support
   - **Health Checks**: Monitors service status and health

2. **Domains** - Self-contained modules that bundle repository, service, and HTTP handlers:
   - Implement `Domain` interface with `GetRepository()` method
   - Example: `auth/domain.go` packages repository, service, handler group, and middleware together
   - Registered via `app.RegisterDomain(name, dbName, domain)`

3. **HTTP Server** (`httpserver/httpserver.go`):
   - Wraps Go's standard HTTP server with middleware support
   - **HandlerGroup**: Composable routing groups with nested middleware
   - Built-in middleware: TraceID, Recovery, Authentication
   - Graceful shutdown with configurable timeout

4. **Database** (`database/database.go`):
   - PostgreSQL via `sqlx` and `lib/pq` driver
   - Automatic migrations through repository registration
   - Repositories implement `migrator` interface with `Migrations()` method
   - Migration tracking ensures idempotent migrations

5. **Queue Processor** (`queue/processor.go`):
   - Generic job processing with worker pool pattern
   - Supports custom queue providers (e.g., `chanqueue.go` for channel-based queues)
   - Graceful shutdown with job draining during shutdown timeout
   - Panic recovery per worker

6. **Scheduler** (`scheduler/scheduler.go`):
   - Executes `application.Runner` implementations at fixed intervals
   - Uses standard `time.Ticker` for scheduling

### Package Structure

- **application/**: Core application lifecycle management
- **auth/**: Authentication domain (session-based, username/password)
- **session/**: Session storage domain
- **httpserver/**: HTTP server with middleware and routing
- **database/**: Database connection and migration system
- **queue/**: Generic job queue processing
- **scheduler/**: Periodic task scheduler
- **openapiserver/**: OpenAPI specification generation and Scalar UI
- **log/**: Logging utilities with context key support
- **demo-app/cmd/**: Example applications showing framework usage
- **internal/cli/**: CLI command implementations

### Key Patterns

**Service Registration and Startup Flow**:
```go
app := application.New()
app.RegisterDatabase("main", db)
app.RegisterDomain("auth", "main", authDomain)
app.RegisterService("api", httpServer)
app.OnStart(migrationTask, config)
app.Run(ctx)
```

**Execution Order**:
1. Database migrations for all registered databases
2. Startup tasks (sequential, with optional abort-on-error)
3. Services (concurrent goroutines)
4. Wait for context cancellation (Ctrl+C)
5. Graceful shutdown of all services

**HandlerGroup Composition**:
HandlerGroups can be nested and mounted on other groups or servers, allowing modular route organization with scoped middleware.

**Context Keys** (`log/log.go`):
- `TraceIDKey`: Request tracing ID
- `ServiceNameKey`: Service identifier
- `StartupTaskKey`: Startup task name
- `WorkerIDKey`: Queue worker ID

## Writing Package Documentation

Package documentation lives in `docs/src/content/docs/packages/` as MDX files. Follow these principles when writing or updating documentation.

### Required Structure

Every package doc should include these sections in order:

1. **Front matter** - YAML with `title: packagename`
2. **Imports** - Starlight components: `LinkButton`, `Steps`, optionally `Code`
3. **Introduction** - Single sentence describing what the package provides
4. **Core Components** - Bulleted list of main types/interfaces with one-line descriptions. Note if a type implements `Runner` interface for Application compatibility
5. **Step-by-step guide** - Using `<Steps>` component with 4-7 numbered steps
6. **Using with Application** - How to integrate with the `application` package
7. **Additional sections** - Package-specific features, best practices, or reference tables
8. **Complete example** - Import from `demo-app/cmd/` or provide inline

### Key Qualities

- **Brief introductions**: One sentence explaining purpose, e.g., "The `queue` package provides a generic, concurrent job processing system"
- **Progressive steps**: Start with creation, move through configuration, end with running and expected output
- **Focused code snippets**: Show only what's necessary, no boilerplate imports unless relevant
- **Expected output**: Include terminal/log output examples to help verify correctness
- **Integration examples**: Always demonstrate `app.RegisterService()` pattern where applicable
- **Prefer demo imports**: Use `import importedCode from '../../../../../demo-app/cmd/example/main.go?raw'` over inline code for complete examples

### Example Step Format

```mdx
<Steps>
1. Create a new instance

    ```go
    s := scheduler.New(time.Second, runner)
    ```

    Brief explanation of parameters or behavior.

2. Next step...
</Steps>
```

## Testing

- Tests use standard Go testing with `testcontainers-go` for PostgreSQL integration tests
- Coverage excludes `demo-app` and `docs` directories
- CI runs on Ubuntu and Windows for cross-platform compatibility
