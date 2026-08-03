package auth

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrUnauthorized    = errors.New("unauthorized action")
)

// Actor represents the currently authenticated user and their context.
// Follows the auth-session building block OS guidelines.
type Actor struct {
	UserID      uuid.UUID
	TenantID    *uuid.UUID // Can be nil for superadmins
	Role        string
	RoleID      string // Added for Casbin
	Permissions map[string]struct{}
}

// HasPermission checks if the actor has a specific permission.
func (a Actor) HasPermission(permission string) bool {
	_, ok := a.Permissions[permission]
	return ok
}

// ActorFromEcho extracts the typed Actor from the Echo context.
func ActorFromEcho(c echo.Context) (Actor, error) {
	val := c.Get("actor")
	if val == nil {
		return Actor{}, ErrUnauthenticated
	}

	actor, ok := val.(Actor)
	if !ok {
		return Actor{}, ErrUnauthenticated
	}

	return actor, nil
}
