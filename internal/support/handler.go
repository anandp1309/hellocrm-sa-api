package support

import (
	"context"
	"net/http"

	"hellocrm-superadmin/internal/platform/database/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Queries *db.Queries
}

func (h *Handler) CreateSupportTicket(c echo.Context) error {
	var req struct {
		TenantID    string `json:"tenant_id"`
		UserID      string `json:"user_id"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	var tenantUUID pgtype.UUID
	tenantUUID.Scan(req.TenantID)

	var userUUID pgtype.UUID
	userUUID.Scan(req.UserID)

	ticket, err := h.Queries.CreateSupportTicket(context.Background(), db.CreateSupportTicketParams{
		TenantID:    tenantUUID,
		UserID:      userUUID,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}

func (h *Handler) ListOpenTickets(c echo.Context) error {
	tickets, err := h.Queries.ListOpenSupportTickets(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tickets,
	})
}

func (h *Handler) ListClosedTickets(c echo.Context) error {
	tickets, err := h.Queries.ListClosedSupportTickets(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tickets,
	})
}

func (h *Handler) ReopenTicket(c echo.Context) error {
	ticketID := c.Param("id")
	var tUUID pgtype.UUID
	tUUID.Scan(ticketID)

	ticket, err := h.Queries.ReopenSupportTicket(context.Background(), tUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Ticket reopened successfully",
		"data":    ticket,
	})
}

func (h *Handler) UpdateStatus(c echo.Context) error {
	ticketID := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	var tUUID pgtype.UUID
	tUUID.Scan(ticketID)

	ticket, err := h.Queries.UpdateSupportTicketStatus(context.Background(), db.UpdateSupportTicketStatusParams{
		ID:     tUUID,
		Status: req.Status,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Ticket status updated successfully",
		"data":    ticket,
	})
}

func (h *Handler) RateTicket(c echo.Context) error {
	ticketID := c.Param("id")
	var req struct {
		CustomerSatisfaction string `json:"customer_satisfaction"`
		Rating               int32  `json:"rating"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	var tUUID pgtype.UUID
	tUUID.Scan(ticketID)

	// Since rating/satisfaction are nullable pgtype, we need to create them
	var pgSatisfaction pgtype.Text
	pgSatisfaction.Scan(req.CustomerSatisfaction)

	var pgRating pgtype.Int4
	pgRating.Scan(req.Rating)

	ticket, err := h.Queries.RateSupportTicket(context.Background(), db.RateSupportTicketParams{
		ID:                   tUUID,
		CustomerSatisfaction: pgSatisfaction,
		Rating:               pgRating,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Ticket rated successfully",
		"data":    ticket,
	})
}

func (h *Handler) GetTicket(c echo.Context) error {
	ticketID := c.Param("id")
	var tUUID pgtype.UUID
	tUUID.Scan(ticketID)

	ticket, err := h.Queries.GetSupportTicket(context.Background(), tUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": ticket,
	})
}

func RegisterRoutes(e *echo.Group, queries *db.Queries) {
	handler := &Handler{Queries: queries}

	supportGroup := e.Group("/support-tickets")
	supportGroup.POST("", handler.CreateSupportTicket)
	supportGroup.GET("/open", handler.ListOpenTickets)
	supportGroup.GET("/closed", handler.ListClosedTickets)
	supportGroup.GET("/:id", handler.GetTicket)
	supportGroup.POST("/:id/reopen", handler.ReopenTicket)
	supportGroup.PUT("/:id/status", handler.UpdateStatus)
	supportGroup.POST("/:id/rate", handler.RateTicket)
}
