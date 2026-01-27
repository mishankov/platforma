# PROJECT KNOWLEDGE BASE

**Generated:** 2026-01-27
**Commit:** e964d40
**Branch:** main

## OVERVIEW

Platforma - Go web application framework providing domain-driven architecture with lifecycle management, HTTP server, database migrations, job queues, and scheduled tasks. Module: `github.com/platforma-dev/platforma`

## STRUCTURE

```
platforma/
├── application/     # Core lifecycle: Application, Domain interface, Runner, Healthcheck
├── auth/            # Auth domain: user CRUD, sessions, middleware, handlers (see auth/AGENTS.md)
├── session/         # Session domain: repository + service pattern
├── httpserver/      # HTTP server with middleware chain, HandlerGroup, graceful shutdown
├── database/        # PostgreSQL with sqlx, migration system, repository registration
├── queue/           # Generic job processor with Handler/Provider interfaces, ChanQueue
├── scheduler/       # Periodic task runner
├── log/             # Structured logging with context keys (traceId, userId, etc.)
├── internal/cli/    # CLI: `platforma generate domain <name>` creates domain scaffolding
├── demo-app/cmd/    # Example apps: api, auth, database, queue, scheduler, clock
└── docs/            # Astro-based documentation site
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new domain | `internal/cli/templates/domain/` | Use `go run . generate domain <name>` |
| Implement Domain interface | `application/domain.go` | Must return `GetRepository() any` |
| Add HTTP routes | `httpserver/handlergroup.go` | `Handle("METHOD /path", handler)` |
| Add middleware | `httpserver/middleware.go` | Implement `Wrap(http.Handler) http.Handler` |
| Database migrations | Repository's `Migrations()` method | Returns `[]database.Migration` |
| Register with app | `application/application.go` | `RegisterDomain()`, `RegisterService()`, `RegisterDatabase()` |
| Background jobs | `queue/processor.go` | Implement `Handler[T]` interface |
| Scheduled tasks | `scheduler/scheduler.go` | Pass `Runner` + duration |
| Logging with context | `log/log.go` | Use `InfoContext(ctx, msg, args...)` |

## CODE MAP

| Package | Key Types | Role |
|---------|-----------|------|
| application | `Application`, `Domain`, `Runner`, `Healthchecker` | Lifecycle orchestration |
| httpserver | `HTTPServer`, `HandlerGroup`, `Middleware` | HTTP layer |
| database | `Database`, `Migration` | PostgreSQL + migrations |
| auth | `Domain`, `Service`, `Repository`, `User`, `AuthenticationMiddleware` | Full auth implementation |
| session | `Domain`, `Service`, `Repository`, `Session` | Session management |
| queue | `Processor[T]`, `Handler[T]`, `Provider[T]`, `ChanQueue[T]` | Job processing |
| scheduler | `Scheduler` | Periodic execution |
| log | `Logger`, context keys | Structured logging |

## CONVENTIONS

- **JSON tags**: camelCase (enforced by tagliatelle linter)
- **Package tests**: Use `_test` suffix packages (e.g., `package auth_test`)
- **Mocks**: Hand-rolled per test file, no external mocking libraries
- **Dependencies**: Interface-based injection, defined locally in each package
- **Domains**: Aggregate Repository + Service (+ optional HandleGroup + Middleware)
- **Constructors**: `New()` or `NewXxx()` functions return structs
- **Error wrapping**: Always wrap with `fmt.Errorf("context: %w", err)`

## ANTI-PATTERNS (THIS PROJECT)

- **No testify assertions** - use standard library comparisons
- **No global state** - except `log.Logger` (linter-exempted)
- **No init functions** - `gochecknoinits` enforced
- **demo-app/ excluded** - not production code, not linted

## LINTING

45+ linters enabled. Key strict rules:
- `err113`: wrap errors, don't create inline
- `exhaustive`: switch must handle all enum cases
- `forcetypeassert`: no unguarded type assertions
- `wrapcheck`: errors from external packages must be wrapped
- `tagliatelle`: JSON tags must be camelCase

## COMMANDS

```bash
# Development (requires task: https://taskfile.dev)
task lint          # Run golangci-lint
task test          # Run tests with coverage HTML
task check         # Lint + test sequentially
task docs          # Start docs dev server (in docs/)

# Code generation
go run . generate domain <name>    # Scaffold new domain in internal/<name>/

# Build demo apps
go build ./demo-app/cmd/api
go build ./demo-app/cmd/auth
# etc.
```

## TESTING

- **Unit tests**: Standard Go with subtests (`t.Run`)
- **Integration tests**: testcontainers-go for PostgreSQL
- **Coverage**: Excludes demo-app/ and docs/
- **Parallel**: Tests use `t.Parallel()` at test and subtest levels

## NOTES

- Go 1.25.0 required
- PostgreSQL driver: `lib/pq` (not pgx)
- No pkg/ directory - framework packages at root for public API
- `internal/cli/` is the only private code
- Documentation in `docs/src/content/docs/packages/*.mdx`
