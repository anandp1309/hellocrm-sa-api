package identity

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"hellocrm-superadmin/internal/auth"
)

type Handler struct {
	userCmd   *UserCommandService
	userQuery *UserQueryService
}

func NewHandler(userCmd *UserCommandService, userQuery *UserQueryService) *Handler {
	return &Handler{
		userCmd:   userCmd,
		userQuery: userQuery,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group, authMW *auth.Middleware) {
	users := g.Group("/users")
	
	// Apply Auth Middleware to all /users endpoints
	users.Use(authMW.RequireAuth)

	// Route level permission bindings based on architecture guidelines
	users.POST("", h.CreateUser, auth.RequirePermission("users", "write"))
	users.GET("/:id", h.GetUser, auth.RequirePermission("users", "read"))
}

type createUserRequest struct {
	Email    string     `json:"email"`
	Password string     `json:"password"`
	RoleID   string     `json:"role_id"`
	TenantID *string    `json:"tenant_id,omitempty"`
}

func (h *Handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	actor, err := auth.ActorFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Role ID format")
	}

	var tenantID *uuid.UUID
	if req.TenantID != nil {
		tid, err := uuid.Parse(*req.TenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Tenant ID format")
		}
		tenantID = &tid
	}

	id, err := h.userCmd.CreateUser(ctx, actor, CreateUserParams{
		Email:    req.Email,
		Password: req.Password,
		RoleID:   roleID,
		TenantID: tenantID,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{"id": id.String()},
	})
}

func (h *Handler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	actor, err := auth.ActorFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User ID format")
	}

	user, err := h.userQuery.GetUserByID(ctx, actor, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}
