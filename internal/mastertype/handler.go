package mastertype

import (
	"net/http"
	"strconv"

	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	masterTypeGroup := g.Group("/master-types")
	
	masterTypeGroup.POST("", h.Create)
	masterTypeGroup.GET("", h.List)
	masterTypeGroup.GET("/:id", h.Get)
	masterTypeGroup.PATCH("/:id", h.Update)
	masterTypeGroup.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateMasterTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	// Ideally you would extract user ID from JWT token in context
	// req.CreatedBy = c.Get("user_id").(string)

	mt, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, mt)
}

func (h *Handler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	mt, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		// Could check for pgx.ErrNoRows here for 404
		return c.JSON(http.StatusNotFound, map[string]string{"error": "master type not found"})
	}

	return c.JSON(http.StatusOK, mt)
}

func (h *Handler) List(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := int32(100)
	offset := int32(0)

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = int32(l)
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = int32(o)
	}

	types, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if types == nil {
		types = []db.MstMasterType{} // return empty array instead of null
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": types,
		"meta": map[string]int32{
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	var req UpdateMasterTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	mt, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mt)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	// In a real scenario, extract from JWT token context
	deletedBy := "" // c.Get("user_id").(string)

	err = h.service.Delete(c.Request().Context(), id, deletedBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
