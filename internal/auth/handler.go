package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	UserID      string   `json:"user_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	token, actor, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	var permissions []string
	for p := range actor.Permissions {
		permissions = append(permissions, p)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token:       token,
		UserID:      actor.UserID.String(),
		Role:        actor.Role,
		Permissions: permissions,
	})
}

func (h *Handler) GetMe(c echo.Context) error {
	actor, err := ActorFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to get user context")
	}

	// Convert the map to a slice of strings for the JSON response
	var permissions []string
	for p := range actor.Permissions {
		permissions = append(permissions, p)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":     actor.UserID,
		"tenant_id":   actor.TenantID,
		"role":        actor.Role,
		"permissions": permissions,
	})
}

// RegisterRoutes registers auth related endpoints.
func (h *Handler) RegisterRoutes(g *echo.Group, authMW *Middleware) {
	authGroup := g.Group("/auth")
	authGroup.POST("/login", h.Login)
	authGroup.GET("/me", h.GetMe, authMW.RequireAuth)
}
