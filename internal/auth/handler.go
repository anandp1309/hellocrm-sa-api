package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hellocrm-superadmin/internal/role"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	UserID      string   `json:"user_id"`
	Role        string   `json:"role"`
	RoleID      string   `json:"role_id"`
	TenantID    *string  `json:"tenant_id,omitempty"`
	TenantName  *string  `json:"tenant_name,omitempty"`
}

type Handler struct {
	authService *Service
	roleQuery   role.QueryService
}

func NewHandler(authService *Service, roleQuery role.QueryService) *Handler {
	return &Handler{
		authService: authService,
		roleQuery:   roleQuery,
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

	// Permissions are intentionally omitted here.
	// The frontend will call /api/v1/auth/me to fetch them.

	var tID *string
	if actor.TenantID != nil {
		idStr := actor.TenantID.String()
		tID = &idStr
	}

	var tName *string
	if actor.TenantName != "" {
		tName = &actor.TenantName
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token:       token,
		UserID:      actor.UserID.String(),
		Role:        actor.Role,
		RoleID:      actor.RoleID,
		TenantID:    tID,
		TenantName:  tName,
	})
}

func (h *Handler) GetMe(c echo.Context) error {
	actor, err := ActorFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to get user context")
	}

	var tID *string
	if actor.TenantID != nil {
		idStr := actor.TenantID.String()
		tID = &idStr
	}

	var tName *string
	if actor.TenantName != "" {
		tName = &actor.TenantName
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":     actor.UserID,
		"tenant_id":   tID,
		"tenant_name": tName,
		"role":        actor.Role,
		"role_id":     actor.RoleID,
	})
}

// RegisterRoutes registers auth related endpoints.
func (h *Handler) RegisterRoutes(g *echo.Group, authMW *Middleware) {
	authGroup := g.Group("/auth")
	authGroup.POST("/login", h.Login)
	authGroup.GET("/me", h.GetMe, authMW.RequireAuth)
}
