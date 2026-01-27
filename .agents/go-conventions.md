# GO CONVENTIONS

Project-specific Go patterns for Platforma.

## JSON TAGS

Use **camelCase** for all JSON struct tags. Enforced by `tagliatelle` linter.

```go
// Correct
type User struct {
    ID        string `json:"id"`
    FirstName string `json:"firstName"`
    CreatedAt time.Time `json:"createdAt"`
}

// Wrong - will fail linting
type User struct {
    FirstName string `json:"first_name"` // snake_case forbidden
}
```

## ERROR HANDLING

### Wrapping Errors

Always wrap errors with context using `fmt.Errorf`:

```go
// Correct
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Wrong - loses context
if err != nil {
    return err
}
```

### No Inline Error Creation (err113)

The `err113` linter forbids creating errors inline in non-test code. Define errors as package-level variables:

```go
// Correct - define at package level
var ErrUserNotFound = errors.New("user not found")

func GetUser(id string) (*User, error) {
    if user == nil {
        return nil, ErrUserNotFound
    }
}

// Wrong - inline error creation forbidden
func GetUser(id string) (*User, error) {
    if user == nil {
        return nil, errors.New("user not found") // err113 violation
    }
}
```

**Note**: Wrapping with `fmt.Errorf("context: %w", err)` is always allowed. Only `errors.New()` inline is forbidden.

## INTERFACES

Define interfaces **locally** in the package that uses them, not where they're implemented:

```go
// In service.go - defines what it needs
type userRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

type Service struct {
    repo userRepository // depends on local interface
}
```

## DEPENDENCIES

Use interface-based dependency injection. Pass dependencies through constructors:

```go
func NewService(repo userRepository, logger *log.Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}
```

## DOMAINS

A domain aggregates related components:
- **Repository** - database operations
- **Service** - business logic
- **HandlerGroup** (optional) - HTTP endpoints
- **Middleware** (optional) - request processing

See `auth/` for a complete example.
