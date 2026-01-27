# LINTING

Strict linting configuration for Platforma. Run with `task lint`.

## KEY RULES

These linters cause the most confusion. Know them.

| Linter | Rule | What It Means |
|--------|------|---------------|
| `err113` | No inline errors | Can't use `errors.New("...")` inline. Define as package vars. |
| `wrapcheck` | Wrap external errors | Errors from other packages must be wrapped with `fmt.Errorf`. |
| `exhaustive` | Exhaustive switches | Switch on enum must handle all cases (or use `default`). |
| `forcetypeassert` | Guarded assertions | No `x.(Type)`. Use `x, ok := y.(Type)` with ok check. |
| `tagliatelle` | JSON camelCase | JSON tags must be camelCase, not snake_case. |
| `gochecknoinits` | No init() | Init functions are forbidden. Use explicit initialization. |
| `gochecknoglobals` | No globals | Global variables forbidden (except `log.Logger`). |

## EXAMPLES

### err113 - Define errors at package level

```go
// Correct
var ErrNotFound = errors.New("not found")

// Wrong
return errors.New("not found")
```

### wrapcheck - Wrap external package errors

```go
// Correct
user, err := repo.GetUser(ctx, id)
if err != nil {
    return fmt.Errorf("get user: %w", err)
}

// Wrong - unwrapped external error
return repo.GetUser(ctx, id)
```

### exhaustive - Handle all enum cases

```go
type Status int
const (
    Active Status = iota
    Inactive
    Deleted
)

// Correct - handles all cases
switch status {
case Active:
    // ...
case Inactive:
    // ...
case Deleted:
    // ...
}

// Also correct - default covers rest
switch status {
case Active:
    // ...
default:
    // ...
}

// Wrong - missing cases
switch status {
case Active:
    // ... 
} // Missing Inactive, Deleted
```

### forcetypeassert - Check type assertions

```go
// Correct
user, ok := ctx.Value(userKey).(*User)
if !ok {
    return nil, ErrNoUser
}

// Wrong - unguarded assertion
user := ctx.Value(userKey).(*User)
```

## EXCLUSIONS

| Path | Excluded Linters | Reason |
|------|------------------|--------|
| `demo-app/` | All | Example code, not production |
| `*_test.go` | `err113`, `errcheck`, `gosec`, etc. | Tests have different requirements |
| `database_test.go` | `paralleltest`, `tparallel` | Integration tests can't always parallelize |

## RUNNING

```bash
task lint       # Check
task fix        # Auto-fix where possible
```
