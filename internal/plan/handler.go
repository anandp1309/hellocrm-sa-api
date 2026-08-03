package plan

import (
	"net/http"
	"strconv"

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
	plans := g.Group("/plans")
	plans.GET("/stats", h.Stats)
	plans.POST("", h.Create)
	plans.GET("/:id", h.Get)
	plans.GET("", h.List)
	plans.PUT("/:id", h.Update)
	plans.DELETE("/:id", h.Delete)
}

func (h *Handler) Stats(c echo.Context) error {
	stats, err := h.queries.GetPlanStats(c.Request().Context())
	if err != nil {
		c.Logger().Error("GetPlanStats error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": stats})
}

type planRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Interval        string  `json:"interval"`
	MaxUsers        int32   `json:"max_users"`
	StorageBytes    int64   `json:"storage_bytes"`
	SmsCredits      int32   `json:"sms_credits"`
	WhatsappCredits int32   `json:"whatsapp_credits"`
	EmailCredits    int32   `json:"email_credits"`
}

func (h *Handler) Create(c echo.Context) error {
	var req planRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	id, err := h.commands.CreatePlan(c.Request().Context(), CreatePlanParams{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Interval:        req.Interval,
		MaxUsers:        req.MaxUsers,
		StorageBytes:    req.StorageBytes,
		SmsCredits:      req.SmsCredits,
		WhatsappCredits: req.WhatsappCredits,
		EmailCredits:    req.EmailCredits,
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

	plan, err := h.queries.GetPlanByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Plan not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": plan})
}

type PlanFilter struct {
	Name         string
	Type         string
	Status       string
	BillingCycle string
	Page         int
	Limit        int
	SortBy       string
	SortOrder    string
}

func (h *Handler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	filter := PlanFilter{
		Name:         c.QueryParam("name"),
		Type:         c.QueryParam("type"),
		Status:       c.QueryParam("status"),
		BillingCycle: c.QueryParam("billing_cycle"),
		Page:         page,
		Limit:        limit,
		SortBy:       c.QueryParam("sort_by"),
		SortOrder:    c.QueryParam("sort_order"),
	}
	if search := c.QueryParam("search"); search != "" {
		filter.Name = search
	}

	paginatedResult, err := h.queries.ListPlans(c.Request().Context(), filter)
	if err != nil {
		c.Logger().Error("ListPlans error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch plans")
	}

	return c.JSON(http.StatusOK, paginatedResult)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	var req planRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.commands.UpdatePlan(c.Request().Context(), UpdatePlanParams{
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Interval:        req.Interval,
		MaxUsers:        req.MaxUsers,
		StorageBytes:    req.StorageBytes,
		SmsCredits:      req.SmsCredits,
		WhatsappCredits: req.WhatsappCredits,
		EmailCredits:    req.EmailCredits,
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

	err = h.commands.DeletePlan(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
