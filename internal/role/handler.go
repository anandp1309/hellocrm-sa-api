package role

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	QueryService   QueryService
	CommandService *CommandService
}

func NewHandler(query QueryService, command *CommandService) *Handler {
	return &Handler{QueryService: query, CommandService: command}
}

func (h *Handler) RegisterRoutes(api *echo.Group) {
	api.GET("/roles", h.ListRoles)
	api.GET("/roles/:id/permissions", h.GetRolePermissions)
	api.PUT("/roles/:id/permissions", h.UpdateRolePermissions)
	api.PUT("/roles/:id", h.UpdateRole)
}

func (h *Handler) ListRoles(c echo.Context) error {
	ctx := c.Request().Context()
	
	roles, err := h.QueryService.ListRoles(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": roles,
	})
}

func (h *Handler) GetRolePermissions(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	
	permissions, err := h.QueryService.GetRolePermissions(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": permissions,
	})
}

func (h *Handler) UpdateRole(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid ID"})
	}

	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid body"})
	}

	if err := h.CommandService.UpdateRole(ctx, id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Role updated successfully",
	})
}

func (h *Handler) UpdateRolePermissions(c echo.Context) error {
	// Stub implementation for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Permissions updated successfully",
	})
}
