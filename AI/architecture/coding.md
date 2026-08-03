# Coding and Backend Architecture

Status: Required

## Dependency direction

The normal backend flow is:

```text
HTTP Handler -> Service -> Repository -> SQLC -> PostgreSQL
```

### Handler responsibilities

Handlers may:

- Parse path, query, header, and body input
- Validate transport-level input
- Obtain the authenticated actor from the trusted authentication boundary
- Call services
- Map application results to HTTP responses

Handlers must not:

- Contain business decisions
- Execute SQL
- Start unrelated background work
- Decide permission logic beyond invoking the approved authorization boundary
- Construct ad-hoc database transactions

### Service responsibilities

Services:

- Own business rules
- Enforce domain invariants
- Enforce action-level authorization
- Own transaction boundaries for multi-step operations
- Coordinate repositories and platform services
- Return typed application errors

### Repository responsibilities

Repositories:

- Encapsulate SQLC calls
- Accept explicit tenant/project scope
- Return domain-neutral persistence results
- Avoid business workflow decisions
- Avoid HTTP concerns

## Package design

- Organise by feature or business capability, not by global technical layer alone.
- Prefer concrete types.
- Add interfaces only for external boundaries, multiple implementations, or valuable test seams.
- Avoid a generic repository abstraction over all tables.
- Avoid “base service” inheritance-style patterns.
- Keep package APIs small.
- Do not create circular dependencies.

## Dependency injection

Use constructor injection:

```go
type UserService struct {
    users UserRepository
    audit AuditWriter
}

func NewUserService(users UserRepository, audit AuditWriter) *UserService {
    return &UserService{
        users: users,
        audit: audit,
    }
}
```

## Context

- Pass `context.Context` as the first parameter of request-scoped functions.
- Do not store context in structs.
- Respect cancellation and deadlines.
- Set timeouts for database and third-party operations.
- Context may carry correlation metadata, but important business inputs such as actor and tenant scope should be explicit service parameters.

## Errors

- Use typed or sentinel application errors.
- Wrap unexpected infrastructure errors with useful operation context.
- Preserve error identity using `%w`.
- Do not wrap repeatedly without adding useful information.
- Do not expose wrapped internal text to clients.
- Log errors once at the application boundary.

Example:

```go
user, err := r.queries.GetUser(ctx, params)
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return User{}, app.ErrNotFound
    }
    return User{}, fmt.Errorf("get user %s: %w", id, err)
}
```

## Concurrency

- Every goroutine must have clear ownership, cancellation, and error handling.
- Do not start uncontrolled goroutines from request handlers.
- Prefer background-job infrastructure for durable work.
- Use database constraints and transactions to protect shared business state.
- Do not rely on frontend state to prevent conflicting updates.

## Testing

- Use table-driven tests for business-rule combinations.
- Test permission failures.
- Test transaction rollback.
- Test concurrency-sensitive paths.
- Prefer integration tests for important SQL behaviour.
