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
	api.POST("/roles", h.CreateRole)
	api.GET("/roles/permissions-template", h.GetPermissionsTemplate)
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

func (h *Handler) CreateRole(c echo.Context) error {
	ctx := c.Request().Context()

	// In a real app, tenantID comes from context/token. Using a dummy for now, or fetching from first tenant.
	// We'll just generate a dummy tenant ID for this SA context
	tenantID := uuid.MustParse("019fd07f-f6d6-776a-8a62-c2789c37b015") 
	
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid body"})
	}

	if err := h.CommandService.CreateRole(ctx, tenantID, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Role created successfully",
	})
}

func (h *Handler) UpdateRolePermissions(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req UpdateRolePermissionsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid body"})
	}

	// Try to get ID from URL first, fallback to payload
	idStr := c.Param("id")
	if idStr == "" || idStr == "undefined" {
		idStr = req.RoleId
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid ID"})
	}

	tenantID := uuid.MustParse("019fd07f-f6d6-776a-8a62-c2789c37b015") 

	if err := h.CommandService.UpdateRolePermissions(ctx, tenantID, id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Permissions updated successfully",
	})
}

func (h *Handler) GetPermissionsTemplate(c echo.Context) error {
	ctx := c.Request().Context()
	
	template, err := h.QueryService.GetPermissionsTemplate(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": template,
	})
}
