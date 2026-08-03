# Authentication Session Building Block

## Security boundary

Authentication middleware validates the session/token and constructs a trusted actor.

```go
type Actor struct {
    UserID      uuid.UUID
    TenantID    uuid.UUID
    ProjectID   *uuid.UUID
    Permissions map[string]struct{}
}
```

The handler retrieves the actor from the trusted request context and passes it explicitly to the service:

```go
func (h *Handler) Update(c echo.Context) error {
    ctx := c.Request().Context()

    actor, err := auth.ActorFromEcho(c)
    if err != nil {
        return h.errors.Write(c, err)
    }

    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return h.errors.Write(c, app.ErrValidation)
    }

    var input UpdateUserInput
    if err := c.Bind(&input); err != nil {
        return h.errors.Write(c, app.ErrValidation)
    }

    result, err := h.service.Update(ctx, actor, id, input)
    if err != nil {
        return h.errors.Write(c, err)
    }

    return c.JSON(http.StatusOK, api.Single(result))
}
```

Service signature:

```go
func (s *Service) Update(
    ctx context.Context,
    actor auth.Actor,
    targetID uuid.UUID,
    input UpdateUserInput,
) (User, error)
```

## Rules

- Never accept user ID or tenant ID from the client as authenticated identity.
- Never let the frontend decide authorization.
- Pass actor/scope explicitly to services.
- Context may still carry request IDs and logging metadata.
- Refresh-token or session revocation must be supported.
- Login, OTP, and password reset endpoints must be rate-limited.
