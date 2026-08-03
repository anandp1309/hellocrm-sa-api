package tenant

import (
	"context"

	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// QueryService handles read operations for tenants.
// This separates read models and logic from the write models (CommandService).
type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{
		repo: repo,
	}
}

type TenantView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *QueryService) GetTenantByID(ctx context.Context, id uuid.UUID) (TenantView, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id.String()); err != nil {
		return TenantView{}, err
	}

	row, err := s.repo.GetTenantByID(ctx, pgID)
	if err != nil {
		return TenantView{}, err
	}

	rowID, _ := uuid.Parse(fmt.Sprintf("%x-%x-%x-%x-%x", row.TenantUuid.Bytes[0:4], row.TenantUuid.Bytes[4:6], row.TenantUuid.Bytes[6:8], row.TenantUuid.Bytes[8:10], row.TenantUuid.Bytes[10:16]))

	return TenantView{
		ID:   rowID,
		Name: row.TenantName,
	}, nil
}



type PaginatedTenants struct {
	Data       []TenantListView `json:"data"`
	Total      int64            `json:"total"`
	Page       int32            `json:"page"`
	Limit      int32            `json:"limit"`
	TotalPages int64            `json:"totalPages"`
}

type TenantListView struct {
	ID           string  `json:"id"`
	Customer     string  `json:"customer"`
	Email        string  `json:"email"`
	Status       string  `json:"status"`
	PlanType     string  `json:"type"`
	Plan         string  `json:"plan"`
	BillingCycle string  `json:"billingCycle"`
	StartDate    string  `json:"startDate"`
	NextRenewal  string  `json:"nextRenewal"`
	MRR          float64 `json:"mrr"`
}

func (s *QueryService) ListTenantsPaginated(ctx context.Context, search, status, planType, plan, billingCycle string, page, limit int32) (PaginatedTenants, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := s.repo.ListTenantsPaginated(ctx, search, status, planType, plan, billingCycle, limit, offset)
	if err != nil {
		return PaginatedTenants{}, err
	}

	var results []TenantListView
	var total int64 = 0

	for _, row := range rows {
		total = row.TotalRecords
		rowID, _ := uuid.Parse(fmt.Sprintf("%x-%x-%x-%x-%x", row.TenantUuid.Bytes[0:4], row.TenantUuid.Bytes[4:6], row.TenantUuid.Bytes[6:8], row.TenantUuid.Bytes[8:10], row.TenantUuid.Bytes[10:16]))
		
		startDate := ""
		if row.StartDate.Valid {
			startDate = row.StartDate.Time.Format("02 Jan 2006")
		}
		
		nextRenewal := "—"
		if row.NextRenewal.Valid {
			nextRenewal = row.NextRenewal.Time.Format("02 Jan 2006")
		}

		mrr := 0.0
		if row.Mrr.Valid {
			f, _ := row.Mrr.Float64Value()
			mrr = f.Float64
		}

		results = append(results, TenantListView{
			ID:           rowID.String(),
			Customer:     row.TenantName,
			Email:        row.EmailAddress.String,
			Status:       row.Status,
			PlanType:     row.PlanType,
			Plan:         row.Plan,
			BillingCycle: row.BillingCycle,
			StartDate:    startDate,
			NextRenewal:  nextRenewal,
			MRR:          mrr,
		})
	}

	totalPages := int64(1)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return PaginatedTenants{
		Data:       results,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

type TenantStats struct {
	TotalCustomers    int64 `json:"totalCustomers"`
	ActiveCustomers   int64   `json:"activeCustomers"`
	TrialCustomers    int64   `json:"trialCustomers"`
	InactiveCustomers int64   `json:"inactiveCustomers"`
	TotalMrr          float64 `json:"totalMrr"`
}

func (s *QueryService) GetTenantStats(ctx context.Context) (TenantStats, error) {
	row, err := s.repo.GetTenantStats(ctx)
	if err != nil {
		return TenantStats{}, err
	}
	var mrr float64
	mrrStr := fmt.Sprintf("%v", row.TotalMrr)
	if parsed, err := strconv.ParseFloat(mrrStr, 64); err == nil {
		mrr = parsed
	}

	return TenantStats{
		TotalCustomers:    row.TotalCustomers,
		ActiveCustomers:   row.ActiveCustomers,
		TrialCustomers:    row.TrialCustomers,
		InactiveCustomers: row.InactiveCustomers,
		TotalMrr:          mrr,
	}, nil
}
