package admin

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	queries  *QueryService
	commands *CommandService
}

func NewHandler(queries *QueryService, commands *CommandService) *Handler {
	return &Handler{
		queries:  queries,
		commands: commands,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	group := g.Group("/admins")
	group.GET("", h.List)
	group.GET("/stats", h.GetStats)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *Handler) List(c echo.Context) error {
	search := c.QueryParam("search")
	role := c.QueryParam("role")
	status := c.QueryParam("status")
	mfa := c.QueryParam("mfa")

	if role == "All Roles" || role == "All" {
		role = ""
	}
	if status == "All Status" || status == "All" {
		status = ""
	}
	if mfa == "All 2FA Status" || mfa == "All" {
		mfa = ""
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	res, err := h.queries.ListPaginated(c.Request().Context(), search, role, status, mfa, int32(page), int32(limit))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list admins")
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStats(c echo.Context) error {
	stats, err := h.queries.GetStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch admin stats")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
		"meta": nil,
	})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	res, err := h.queries.GetAdmin(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Admin not found")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateAdminRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	id, err := h.commands.CreateAdmin(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Admin created successfully",
		"id":      id,
	})
}

func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	var req UpdateAdminRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.commands.UpdateAdmin(c.Request().Context(), id, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Admin updated successfully",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err := h.commands.DeleteAdmin(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Admin deleted successfully",
	})
}
