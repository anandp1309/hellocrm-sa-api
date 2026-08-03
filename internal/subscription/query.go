package subscription

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"hellocrm-superadmin/internal/platform/database/db"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

type Subscription struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	PlanID     string    `json:"plan_id"`
	Status     string    `json:"status"`
}

type SubscriptionStats struct {
	Total     int    `json:"total"`
	Active    int    `json:"active"`
	Trial     int    `json:"trial"`
	Expired   int    `json:"expired"`
	MRR       string `json:"mrr"`
	MRRGrowth string `json:"mrrGrowth"`
}

func (s *QueryService) GetSubscriptionStats(ctx context.Context) (SubscriptionStats, error) {
	// STUB: Replace with actual database call
	return SubscriptionStats{
		Total:     50,
		Active:    45,
		Trial:     3,
		Expired:   2,
		MRR:       "₹ 1,50,000",
		MRRGrowth: "+12%",
	}, nil
}

func (s *QueryService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (Subscription, error) {
	if id == uuid.Nil {
		return Subscription{}, fmt.Errorf("invalid id")
	}
	// STUB: Replace with actual database call
	return Subscription{
		ID:         id,
		CustomerID: uuid.Must(uuid.NewV7()),
		PlanID:     "pro-plan",
		Status:     "active",
	}, nil
}

func (s *QueryService) ListSubscriptions(ctx context.Context, params db.ListSubscriptionsParams) ([]db.ListSubscriptionsRow, error) {
	return s.repo.queries.ListSubscriptions(ctx, params)
}
