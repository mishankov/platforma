# TESTING

Testing patterns and requirements for Platforma.

## TEST PACKAGE NAMING

Use `_test` suffix packages for external testing:

```go
// Correct - in auth/handler_test.go
package auth_test

import "github.com/platforma-dev/platforma/auth"

// Wrong - internal test package
package auth
```

## PARALLEL TESTS

All tests must use `t.Parallel()` at both test and subtest levels:

```go
func TestUserService(t *testing.T) {
    t.Parallel() // Required at test level

    t.Run("creates user", func(t *testing.T) {
        t.Parallel() // Required at subtest level
        // ...
    })
}
```

Enforced by `paralleltest` and `tparallel` linters.

## ASSERTIONS

Use standard library comparisons. **No testify**.

```go
// Correct
if got != want {
    t.Fatalf("expected %v, got %v", want, got)
}

if w.Code != http.StatusOK {
    t.Fatalf("expected status 200, got %d", w.Code)
}

// Wrong - no testify
assert.Equal(t, want, got)
require.NoError(t, err)
```

## MOCKING

Hand-roll mocks per test file. No external mocking libraries.

```go
// Define mock in test file
type mockUserRepository struct {
    getUserErr  error
    getUserUser *auth.User
}

func (m *mockUserRepository) GetUser(ctx context.Context, id string) (*auth.User, error) {
    return m.getUserUser, m.getUserErr
}

// Use in test
func TestService_GetUser(t *testing.T) {
    t.Parallel()

    mock := &mockUserRepository{
        getUserUser: &auth.User{ID: "123"},
    }
    svc := auth.NewService(mock)
    // ...
}
```

## INTEGRATION TESTS

Use `testcontainers-go` for PostgreSQL integration tests:

```go
import "github.com/testcontainers/testcontainers-go/modules/postgres"

func TestDatabaseIntegration(t *testing.T) {
    ctx := context.Background()
    container, err := postgres.Run(ctx, "postgres:16")
    if err != nil {
        t.Fatal(err)
    }
    defer container.Terminate(ctx)

    connStr, _ := container.ConnectionString(ctx)
    // Use connStr for database connection
}
```

## COVERAGE

Coverage excludes:
- `demo-app/` - example code, not production
- `docs/` - documentation site

Run with: `task test` (generates `coverage.html`)
