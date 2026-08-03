package universal

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	universalGroup := g.Group("/universal-masters")
	
	// Create a new universal value linked to a master_type
	universalGroup.POST("", h.Create)
	universalGroup.PUT("/:id", h.Update)
	universalGroup.GET("/master-types/:id", h.GetByMasterType)
	universalGroup.DELETE("/:id", h.Deactivate)
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateUniversalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	mt, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, mt)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	
	var req CreateUniversalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "updated successfully"})
}

func (h *Handler) GetByMasterType(c echo.Context) error {
	id := c.Param("id")
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	filter := UniversalFilter{
		Search:    c.QueryParam("search"),
		Scope:     c.QueryParam("scope"),
		Tenant:    c.QueryParam("tenant"),
		Status:    c.QueryParam("status"),
		Page:      page,
		Limit:     limit,
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}
	
	paginatedResult, err := h.service.GetByMasterType(c.Request().Context(), id, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, paginatedResult)
}

func (h *Handler) Deactivate(c echo.Context) error {
	id := c.Param("id")
	err := h.service.Deactivate(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "deactivated successfully"})
}
