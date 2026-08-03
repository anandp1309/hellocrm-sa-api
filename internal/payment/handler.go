package payment

import (
	"net/http"
	"github.com/google/uuid"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	queries *QueryService
}

func NewHandler(queries *QueryService) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	payments := g.Group("/payments")
	payments.GET("", h.List)
	payments.GET("/stats", h.GetStats)
	payments.POST("", h.Create)
	payments.GET("/:id", h.Get)
	payments.PUT("/:id", h.Update)
	payments.DELETE("/:id", h.Delete)
}

func (h *Handler) List(c echo.Context) error {
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	method := c.QueryParam("method")
	plan := c.QueryParam("plan")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	if status == "All Status" {
		status = ""
	}
	if method == "All Methods" {
		method = ""
	}
	if plan == "All Plans" {
		plan = ""
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	res, err := h.queries.ListPaymentsPaginated(c.Request().Context(), search, status, method, plan, startDate, endDate, int32(page), int32(limit))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list payments")
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStats(c echo.Context) error {
	stats, err := h.queries.GetPaymentStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch payment stats")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
		"meta": nil,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	
	payment, err := h.queries.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusCreated, payment)
}

func (h *Handler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}
	
	payment, err := h.queries.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "payment not found"})
	}
	
	return c.JSON(http.StatusOK, payment)
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}
	
	var req UpdatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	
	payment, err := h.queries.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusOK, payment)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}
	
	// Assume user uuid could be taken from context, but we will leave it empty for now
	if err := h.queries.Delete(c.Request().Context(), id, ""); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.NoContent(http.StatusNoContent)
}
