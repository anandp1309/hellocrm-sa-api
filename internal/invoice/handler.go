package invoice

import (
	"net/http"
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
	group := g.Group("/invoices")
	group.GET("", h.List)
	group.GET("/stats", h.GetStats)
	group.GET("/:id", h.Get)
}

func (h *Handler) List(c echo.Context) error {
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	if status == "All Status" || status == "All" {
		status = ""
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	res, err := h.queries.ListPaginated(c.Request().Context(), search, status, startDate, endDate, int32(page), int32(limit))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list invoices")
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStats(c echo.Context) error {
	stats, err := h.queries.GetStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch invoice stats")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
		"meta": nil,
	})
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	res, err := h.queries.Get(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invoice not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
}
