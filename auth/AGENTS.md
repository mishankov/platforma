# AUTH DOMAIN

Full authentication domain with user management, session handling, and HTTP middleware.

## STRUCTURE

```
auth/
├── domain.go      # Domain aggregate: Repository + Service + HandlerGroup + Middleware
├── model.go       # User struct, Status enum (Active/Inactive/Deleted)
├── repository.go  # DB operations, migrations (users table)
├── service.go     # Business logic: register, login, logout, password change
├── middleware.go  # AuthenticationMiddleware - validates session, injects user to context
├── handler_*.go   # HTTP handlers: register, login, logout, get, change_password, delete
├── context.go     # Context helpers: UserFromContext, SetUserToContext
├── errors.go      # Domain errors: ErrUserNotFound, ErrInvalidCredentials, etc.
└── cleanup.go     # Session cleanup job for queue processing
```

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Add auth endpoint | `domain.go` | Add to `authAPI.Handle()` calls |
| Change user model | `model.go` + `repository.go` | Update struct + migrations |
| Modify auth logic | `service.go` | All business rules here |
| Protect routes | Use `domain.Middleware` | Wraps handlers, requires valid session |
| Access current user | `context.go` | `auth.UserFromContext(ctx)` |
| Custom validation | `service.go` | Pass validators to `New()` constructor |

## INTERFACES (Dependencies)

```go
// Required by Repository
type db interface {
    QueryRowContext(ctx, query, args...) *sql.Row
    ExecContext(ctx, query, args...) (sql.Result, error)
}

// Required by Service  
type authStorage interface {
    CreateSession(ctx, userID) (*session.Session, error)
    DeleteSession(ctx, sessionID) error
    GetSession(ctx, sessionID) (*session.Session, error)
}
```

## MIDDLEWARE USAGE

```go
// Protect entire handler group
protectedAPI := httpserver.NewHandlerGroup()
protectedAPI.Use(authDomain.Middleware)

// Protect single handler
protectedHandler := authDomain.Middleware.Wrap(myHandler)
```

## SESSION COOKIE

Cookie name configurable via `New()` constructor. Default flow:
1. Login → creates session → sets cookie
2. Middleware reads cookie → validates session → injects user to context
3. Logout → deletes session → clears cookie

## CLEANUP JOB

`CleanupJob` implements queue handler for expired session cleanup. Enqueue via `cleanupEnqueuer` interface passed to `New()`.
