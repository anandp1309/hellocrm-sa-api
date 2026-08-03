package usage

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	QueryService QueryService
}

func NewHandler(query QueryService) *Handler {
	return &Handler{QueryService: query}
}

func (h *Handler) RegisterRoutes(api *echo.Group) {
	api.GET("/usage", h.GetUsage)
}

func (h *Handler) GetUsage(c echo.Context) error {
	stats, err := h.QueryService.GetUsageStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}
