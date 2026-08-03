package tenant

import (
	"fmt"
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
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	tenants := g.Group("/tenants")
	tenants.POST("", h.Create)
	tenants.GET("", h.List)
	tenants.GET("/stats", h.GetStats)
	tenants.GET("/:id", h.Get)
	tenants.PATCH("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}

type createTenantRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	ContactPersonName string `json:"contactPersonName"`
	MobileNumber      string `json:"mobileNumber"`
	CountryName       string `json:"countryName"`
	StateName         string `json:"stateName"`
	CityName          string `json:"cityName"`
	Address           string `json:"address"`
	GstNumber         string `json:"gstNumber"`
	Remarks           string `json:"remarks"`
}

func (h *Handler) Create(c echo.Context) error {
	var req createTenantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

		id, err := h.commands.CreateTenant(c.Request().Context(), CreateTenantParams{
		Name:              req.Name,
		Email:             req.Email,
		ContactPersonName: req.ContactPersonName,
		MobileNumber:      req.MobileNumber,
		CountryName:       req.CountryName,
		StateName:         req.StateName,
		CityName:          req.CityName,
		Address:           req.Address,
		GstNumber:         req.GstNumber,
		Remarks:           req.Remarks,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"id": id.String(),
		},
		"meta": nil,
	})
}

func (h *Handler) List(c echo.Context) error {
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	planType := c.QueryParam("planType")
	plan := c.QueryParam("plan")
	billingCycle := c.QueryParam("billingCycle")

	if status == "All Status" {
		status = ""
	}
	if planType == "All" || planType == "All Type" || planType == "All Types" {
		planType = ""
	}
	if plan == "All" || plan == "All Plan" || plan == "All Plans" {
		plan = ""
	}
	if billingCycle == "All" || billingCycle == "All Cycle" {
		billingCycle = ""
	}

	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	var page, limit int32 = 1, 10
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = int32(p)
		}
	}
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = int32(l)
		}
	}

	res, err := h.queries.ListTenantsPaginated(c.Request().Context(), search, status, planType, plan, billingCycle, page, limit)
	if err != nil {
		fmt.Printf("Error in ListTenantsPaginated: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list tenants")
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID format")
	}

	tenant, err := h.queries.GetTenantByID(c.Request().Context(), id)
	if err != nil {
		// Log internal error here
		return echo.NewHTTPError(http.StatusNotFound, "Tenant not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tenant,
		"meta": nil,
	})
}

type updateTenantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID format")
	}

	var req updateTenantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.commands.UpdateTenant(c.Request().Context(), UpdateTenantParams{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update tenant")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tenant updated successfully",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID format")
	}

	err = h.commands.DeleteTenant(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete tenant")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetStats(c echo.Context) error {
	stats, err := h.queries.GetTenantStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": stats,
		"meta": nil,
	})
}
