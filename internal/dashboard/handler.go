package dashboard

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
	api.GET("/dashboard", h.GetDashboard)
}

func (h *Handler) GetDashboard(c echo.Context) error {
	stats, err := h.QueryService.GetDashboardStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}
