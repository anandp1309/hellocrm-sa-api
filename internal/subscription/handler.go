package subscription

import (
	"net/http"
	"strconv"
	"math"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"hellocrm-superadmin/internal/platform/database/db"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(commands *CommandService, queries *QueryService) *Handler {
	return &Handler{commands: commands, queries: queries}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	subs := g.Group("/subscriptions")
	subs.POST("", h.Create)
	subs.GET("/stats", h.Stats)
	subs.GET("/:id", h.Get)
	subs.GET("", h.List)
	subs.PUT("/:id", h.Update)
	subs.DELETE("/:id", h.Cancel) // Often cancel rather than hard delete for subscriptions
	
	// Mock endpoints for Subscription Details Sections
	subs.GET("/:id/details", h.GetDetailsMock)
	subs.GET("/:id/addons", h.GetAddonsMock)
	subs.GET("/:id/renewals", h.GetRenewalsMock)
	subs.GET("/:id/payments", h.GetPaymentsMock)
	subs.GET("/:id/usage", h.GetUsageMock)
	subs.GET("/:id/activity", h.GetActivityMock)
}

type createSubscriptionRequest struct {
	CustomerID string `json:"customer_id"`
	PlanID     string `json:"plan_id"`
	Status     string `json:"status"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	AmountPaid string `json:"amount_paid"`
}

type updateSubscriptionRequest struct {
	PlanID string `json:"plan_id"`
	Status string `json:"status"`
}

func (h *Handler) Create(c echo.Context) error {
	var req createSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer UUID")
	}

	id, err := h.commands.CreateSubscription(c.Request().Context(), CreateSubscriptionParams{
		CustomerID: customerID,
		PlanID:     req.PlanID,
		Status:     req.Status,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		AmountPaid: req.AmountPaid,
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

	sub, err := h.queries.GetSubscriptionByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Subscription not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": sub})
}

func (h *Handler) Stats(c echo.Context) error {
	stats, err := h.queries.GetSubscriptionStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": stats})
}

func (h *Handler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	params := db.ListSubscriptionsParams{
		SearchName:    c.QueryParam("search"),
		Status:        c.QueryParam("status"),
		PlanName:      c.QueryParam("plan"),
		BillingCycle:  c.QueryParam("billing_cycle"),
		PaymentStatus: c.QueryParam("payment_status"),
		Offset:        int32(offset),
		Limit:         int32(limit),
	}

	subs, err := h.queries.ListSubscriptions(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch subscriptions")
	}

	var totalRecords int64 = 0
	if len(subs) > 0 {
		totalRecords = subs[0].TotalRecords
	}
	
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Map DB rows to JSON response
	type SubResponse struct {
		ID              string `json:"id"`
		Customer        string `json:"customer"`
		Email           string `json:"email"`
		SubscriptionNo  string `json:"subscriptionNumber"`
		Plan            string `json:"plan"`
		PlanType        string `json:"planType"`
		BillingCycle    string `json:"billingCycle"`
		Status          string `json:"status"`
		StartDate       string `json:"startDate"`
		NextBillingDate string `json:"nextBillingDate"`
		MRR             string `json:"mrr"`
		PaymentStatus   string `json:"paymentStatus"`
		Addons          int    `json:"addons"`
		Initials        string `json:"initials"`
		IconBg          string `json:"iconBg"`
		IconColor       string `json:"iconColor"`
	}

	var data []SubResponse
	for i, s := range subs {
		// Mock visual fields
		initials := "CO"
		if len(s.CustomerName) >= 2 {
			initials = s.CustomerName[:2]
		}
		
		bgColors := []string{"bg-blue-50", "bg-green-50", "bg-purple-50", "bg-orange-50"}
		textColors := []string{"text-blue-600", "text-green-600", "text-purple-600", "text-orange-600"}
		colorIdx := i % len(bgColors)

		data = append(data, SubResponse{
			ID:              s.SubscriptionNumber,
			Customer:        s.CustomerName,
			Email:           s.CustomerEmail.String,
			SubscriptionNo:  s.SubscriptionNumber,
			Plan:            s.PlanName,
			PlanType:        s.PlanType.String,
			BillingCycle:    s.BillingCycle.String,
			Status:          s.Status.String,
			StartDate:       s.StartDate.Time.Format("2006-01-02"),
			NextBillingDate: s.NextBillingDate.Time.Format("2006-01-02"),
			MRR:             "₹ " + s.Mrr.Int.String(),
			PaymentStatus:   s.PaymentStatus.String,
			Addons:          0, // Stub addons
			Initials:        initials,
			IconBg:          bgColors[colorIdx],
			IconColor:       textColors[colorIdx],
		})
	}
	
	if data == nil {
		data = []SubResponse{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":       data,
		"total":      totalRecords,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
	}

	var req updateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.commands.UpdateSubscription(c.Request().Context(), UpdateSubscriptionParams{
		ID:     id,
		PlanID: req.PlanID,
		Status: req.Status,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Cancel(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing subscription ID")
	}

	err := h.commands.CancelSubscription(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cancel subscription")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Subscription cancelled"})
}

// MOCK ENDPOINTS for UI UI

func (h *Handler) GetDetailsMock(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"planName": "Professional Plan",
			"planType": "Builder",
			"billingCycle": "Monthly",
			"usersIncluded": 10,
			"workspacesIncluded": 5,
			"trialDays": 0,
			"planAmount": "₹ 9,999 / month",
			"companyName": "Green Valley Builders",
			"contactPerson": "Amit Sharma",
			"email": "manager@greenvalley.com",
			"phone": "+91 98765 43210",
			"gstNumber": "27ABCDE1234F1Z5",
			"customerSince": "15 May 2026",
			"totalWorkspaces": 3,
			"totalUsers": 7,
		},
	})
}

func (h *Handler) GetAddonsMock(c echo.Context) error {
	addons := []map[string]interface{}{
		{"name": "Extra Users Pack", "quantity": "5 Users", "cycle": "Monthly", "expiryDate": "30 Jun 2026", "status": "Active"},
		{"name": "SMS Pack - 5000", "quantity": "1 Pack", "cycle": "Monthly", "expiryDate": "30 Jun 2026", "status": "Active"},
		{"name": "Storage 20 GB", "quantity": "1 Pack", "cycle": "Monthly", "expiryDate": "30 Jun 2026", "status": "Active"},
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": addons})
}

func (h *Handler) GetRenewalsMock(c echo.Context) error {
	renewals := []map[string]interface{}{
		{"date": "01 Jul 2026", "plan": "Professional Plan", "cycle": "Monthly", "duration": "01 Jun 2026 - 30 Jun 2026", "amount": "₹ 9,999", "status": "Paid", "renewedBy": "Super Admin", "isUpgraded": false},
		{"date": "01 Jun 2026", "plan": "Professional Plan", "cycle": "Monthly", "duration": "01 May 2026 - 31 May 2026", "amount": "₹ 9,999", "status": "Paid", "renewedBy": "Super Admin", "isUpgraded": false},
		{"date": "15 Feb 2026", "plan": "Upgraded to Professional Plan", "cycle": "Monthly", "duration": "15 Feb 2026 - 28 Feb 2026", "amount": "₹ 5,000", "status": "Paid", "renewedBy": "Super Admin", "isUpgraded": true},
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"totalRenewals": 3,
			"onTime": "3 (100%)",
			"upcomingRenewal": "01 Jul 2026",
			"status": "Active",
			"history": renewals,
		},
	})
}

func (h *Handler) GetPaymentsMock(c echo.Context) error {
	payments := []map[string]interface{}{
		{"invoiceNo": "INV-2026-0612", "invoiceDate": "01 Jul 2026", "amount": "₹ 8,473.73", "gst": "₹ 1,522.27", "totalAmount": "₹ 9,999.00", "paymentMethod": "Razorpay", "paymentDate": "01 Jul 2026", "status": "Paid"},
		{"invoiceNo": "INV-2026-0131", "invoiceDate": "31 Jan 2026", "amount": "₹ 0.00", "gst": "₹ 0.00", "totalAmount": "₹ 0.00", "paymentMethod": "—", "paymentDate": "31 Jan 2026", "status": "Free Trial"},
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": payments})
}

func (h *Handler) GetUsageMock(c echo.Context) error {
	usage := []map[string]interface{}{
		{"type": "Users", "used": 7, "allocated": 10, "percentage": 70},
		{"type": "Workspaces", "used": 3, "allocated": 5, "percentage": 60},
		{"type": "Storage", "used": "12.4 GB", "allocated": "25 GB", "percentage": 50},
		{"type": "SMS", "used": "2,150", "allocated": "5,000", "percentage": 43},
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": usage})
}

func (h *Handler) GetActivityMock(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"data": []string{}})
}
