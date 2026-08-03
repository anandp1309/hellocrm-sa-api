package addon

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(commands *CommandService, queries *QueryService) *Handler {
	return &Handler{commands: commands, queries: queries}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	addons := g.Group("/addons")
	addons.POST("", h.Create)
	addons.GET("/:id", h.Get)
	addons.GET("", h.List)
	addons.PUT("/:id", h.Update)
	addons.DELETE("/:id", h.Delete)
}

type addonRequest struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Limit     string `json:"limit"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	IconSvg   string `json:"icon_svg"`
	IconColor string `json:"icon_color"`
}

func (h *Handler) Create(c echo.Context) error {
	var req addonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	id, err := h.commands.CreateAddon(c.Request().Context(), CreateAddonParams{
		Name:      req.Name,
		Category:  req.Category,
		Limit:     req.Limit,
		Price:     req.Price,
		Status:    req.Status,
		IconSvg:   req.IconSvg,
		IconColor: req.IconColor,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{"id": id.String()},
	})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	addon, err := h.queries.GetAddonByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Addon not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": addon})
}

func (h *Handler) List(c echo.Context) error {
	addons, err := h.queries.ListAddons(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch addons")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": addons})
}

func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	var req addonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	params := UpdateAddonParams{ID: id}
	if req.Name != "" { params.Name = &req.Name }
	if req.Category != "" { params.Category = &req.Category }
	if req.Limit != "" { params.Limit = &req.Limit }
	if req.Price != "" { params.Price = &req.Price }
	if req.Status != "" { params.Status = &req.Status }
	if req.IconSvg != "" { params.IconSvg = &req.IconSvg }
	if req.IconColor != "" { params.IconColor = &req.IconColor }

	err = h.commands.UpdateAddon(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	err = h.commands.DeleteAddon(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
