package subscription

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"hellocrm-superadmin/internal/platform/database/db"
)

type CommandService struct {
	repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{repo: repo}
}

type CreateSubscriptionParams struct {
	CustomerID uuid.UUID
	PlanID     string
	Status     string
	StartDate  string
	EndDate    string
	AmountPaid string
}

type UpdateSubscriptionParams struct {
	ID     uuid.UUID
	PlanID string
	Status string
}

func (s *CommandService) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) (uuid.UUID, error) {
	if params.CustomerID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("customer id is required")
	}
	planUUID, err := uuid.Parse(params.PlanID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid plan id")
	}

	// Fetch plan_price_uuid
	var pUuid pgtype.UUID
	pUuid.Bytes = planUUID
	pUuid.Valid = true

	planPriceUUID, err := s.repo.queries.GetFirstPlanPriceByPlan(ctx, pUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not find plan price: %w", err)
	}

	// Parse dates
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, params.StartDate)
	if err != nil {
		startDate = time.Now()
	}
	endDate, err := time.Parse(layout, params.EndDate)
	if err != nil {
		endDate = startDate.AddDate(0, 1, 0)
	}

	var sd pgtype.Date
	sd.Time = startDate
	sd.Valid = true

	var ed pgtype.Date
	ed.Time = endDate
	ed.Valid = true

	var tUuid pgtype.UUID
	tUuid.Bytes = params.CustomerID
	tUuid.Valid = true

	newId, _ := uuid.NewV7()
	var nUuid pgtype.UUID
	nUuid.Bytes = newId
	nUuid.Valid = true

	subNum := "SUB-" + strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(1000+rand.Intn(9000))

	var amount pgtype.Numeric
	_, err = strconv.ParseFloat(params.AmountPaid, 64)
	if err == nil {
		amount.Int = nil // For exact precision we'd use math/big, but as a shortcut we can pass string to Postgres
		_ = amount.Scan(params.AmountPaid)
	} else {
		amount.Valid = false
	}

	_, err = s.repo.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		TenantSubscriptionUuid: nUuid,
		SubscriptionNumber:     subNum,
		TenantUuid:             tUuid,
		PlanPriceUuid:          planPriceUUID,
		SubscriptionStartDate:  sd,
		SubscriptionEndDate:    ed,
		AmountPaid:             amount,
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return newId, nil
}

func (s *CommandService) UpdateSubscription(ctx context.Context, params UpdateSubscriptionParams) error {
	if params.ID == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	// STUB: Replace with actual database call for updating
	return nil
}

func (s *CommandService) CancelSubscription(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.repo.queries.CancelSubscription(ctx, id)
}
