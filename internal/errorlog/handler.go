package errorlog

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api *echo.Group) {
	group := api.Group("/error-logs")
	group.GET("", h.List)
}

func (h *Handler) List(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": []string{},
	})
}
