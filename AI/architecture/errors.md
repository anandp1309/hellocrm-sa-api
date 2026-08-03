# Error Handling Standard

Status: Required

## Application errors

Use stable machine-readable error codes. Recommended categories:

```go
var (
    ErrNotFound       = errors.New("not found")
    ErrValidation     = errors.New("validation failed")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrConflict       = errors.New("conflict")
    ErrDuplicate      = errors.New("duplicate")
    ErrInvalidState   = errors.New("invalid state")
    ErrUnavailable    = errors.New("service unavailable")
)
```

## Layer behaviour

### Repository

- Convert expected database conditions into known application errors.
- Wrap unexpected database errors.
- Do not produce user-facing messages.

### Service

- Preserve known application errors.
- Add domain context only when useful.
- Enforce business-state conflicts explicitly.
- Do not log errors that will be logged by the application boundary.

### Handler

- Convert application errors to the standard API error contract.
- Assign an HTTP status.
- Include request ID.
- Never expose SQL, stack traces, or internal error strings.

## HTTP mapping

| Application error | HTTP status | API code |
|---|---:|---|
| Validation | 422 | `VALIDATION_ERROR` |
| Unauthorized | 401 | `UNAUTHORIZED` |
| Forbidden | 403 | `FORBIDDEN` |
| Not found | 404 | `NOT_FOUND` |
| Duplicate | 409 | `DUPLICATE_RESOURCE` |
| Conflict | 409 | `CONFLICT` |
| Invalid state | 409 | `INVALID_STATE` |
| Rate limited | 429 | `RATE_LIMITED` |
| Unavailable | 503 | `SERVICE_UNAVAILABLE` |
| Unexpected | 500 | `INTERNAL_ERROR` |

## Validation errors

Field validation should use a stable structure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed.",
    "fields": {
      "email": "Enter a valid email address.",
      "mobile": "Mobile number is required."
    },
    "request_id": "req_01..."
  }
}
```

## Cancellation and timeout

- Map client cancellation separately from server failure in logs.
- Treat `context.DeadlineExceeded` as a timeout.
- Do not retry automatically inside handlers.
- Retry policy belongs in approved integration or background-job code.

## Database constraints

Map known constraints to stable errors. Do not return raw constraint names to clients.
