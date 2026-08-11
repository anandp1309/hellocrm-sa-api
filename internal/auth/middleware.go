package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Middleware provides Echo middleware for extracting tokens and building the Actor.
type Middleware struct {
	authService *Service
	enforcer    *CasbinEnforcer
}

func NewMiddleware(authService *Service, enforcer *CasbinEnforcer) *Middleware {
	return &Middleware{
		authService: authService,
		enforcer:    enforcer,
	}
}

// RequireAuth enforces that a valid token exists and injects the Actor into the context.
func (m *Middleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		userID, roleID, err := m.authService.ParseToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		actor, err := m.authService.GetActorContext(c.Request().Context(), userID, roleID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load permissions")
		}

		// Inject enforcer pointer inside context so RequirePermission can use it
		c.Set("casbin", m.enforcer)
		c.Set("actor", actor)
		return next(c)
	}
}

// RequirePermission is a middleware to check if the actor has a specific permission via Casbin.
func RequirePermission(obj, act string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			actor, err := ActorFromEcho(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			val := c.Get("casbin")
			if val == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "casbin enforcer not found")
			}
			enforcer := val.(*CasbinEnforcer)

			// Super Admin override or standard Casbin check
			// We set sub == "super_admin" as allowed by default in the Casbin model matcher.
			allowed, err := enforcer.Enforce(actor.RoleID, obj, act)
			if err != nil || !allowed {
				return echo.NewHTTPError(http.StatusForbidden, ErrUnauthorized.Error())
			}

			return next(c)
		}
	}
}
