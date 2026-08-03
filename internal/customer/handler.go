package customer

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
	customers := g.Group("/customers")
	customers.POST("", h.Create)
	customers.GET("/:id", h.Get)
	customers.GET("", h.List)
	customers.PUT("/:id", h.Update)
	customers.DELETE("/:id", h.Delete)
}

type customerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) Create(c echo.Context) error {
	var req customerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	id, err := h.commands.CreateCustomer(c.Request().Context(), CreateCustomerParams{
		Name:  req.Name,
		Email: req.Email,
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

	cust, err := h.queries.GetCustomerByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": cust})
}

func (h *Handler) List(c echo.Context) error {
	customers, err := h.queries.ListCustomers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch customers")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}

func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	var req customerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.commands.UpdateCustomer(c.Request().Context(), UpdateCustomerParams{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
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

	err = h.commands.DeleteCustomer(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
