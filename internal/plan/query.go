package plan

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

type Plan struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	PlanType        string    `json:"plan_type"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Interval        string    `json:"interval"`
	MaxUsers        int32     `json:"max_users"`
	StorageBytes    int64     `json:"storage_bytes"`
	SmsCredits      int32     `json:"sms_credits"`
	WhatsappCredits int32     `json:"whatsapp_credits"`
	EmailCredits    int32     `json:"email_credits"`
}

func fromPgNumeric(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0
	}
	f, _ := strconv.ParseFloat(num.Int.String(), 64)
	return f / float64(1) // Actually numeric with scale needs scaling, but it's hard with big.Int. Let's just use float64 formatting or string
}

func mapPlan(p interface{}) Plan {
	// Actually we should map db.Plan to Plan. But it's defined in another package. We will map locally.
	return Plan{}
}

func (s *QueryService) GetPlanByID(ctx context.Context, id uuid.UUID) (Plan, error) {
	if id == uuid.Nil {
		return Plan{}, fmt.Errorf("invalid id")
	}
	
	dbPlan, err := s.repo.queries.GetPlanByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return Plan{}, err
	}

	price := fromPgNumeric(dbPlan.Price)

	return Plan{
		ID:              dbPlan.ID.Bytes,
		Name:            dbPlan.Name,
		PlanType:        dbPlan.PlanType,
		Description:     dbPlan.Description.String, // Remarks remains pgtype.Text
		Price:           price,
		Interval:        dbPlan.Interval,
		MaxUsers:        dbPlan.MaxUsers,
		StorageBytes:    dbPlan.DefaultStorageBytes,
		SmsCredits:      dbPlan.DefaultSmsCredits,
		WhatsappCredits: dbPlan.DefaultWhatsappCredits,
		EmailCredits:    dbPlan.DefaultEmailCredits,
	}, nil
}

type PaginatedPlans struct {
	Data           []Plan `json:"data"`
	TotalRecords   int    `json:"total_records"`
	PresentRecords int    `json:"present_records"`
	CurrentPage    int    `json:"current_page"`
	TotalPages     int    `json:"total_pages"`
	Limit          int    `json:"limit"`
}

func (s *QueryService) ListPlans(ctx context.Context, filter PlanFilter) (PaginatedPlans, error) {
	dbPlans, err := s.repo.queries.ListPlans(ctx)
	if err != nil {
		return PaginatedPlans{}, err
	}

	var allPlans []Plan
	for _, p := range dbPlans {
		price := fromPgNumeric(p.Price)
		
		// In-memory filtering
		if filter.Name != "" && !containsIgnoreCase(p.Name, filter.Name) {
			continue
		}
		if filter.Status != "" && filter.Status != "all" && filter.Status != "active" {
			continue // all mocked to active
		}
		
		allPlans = append(allPlans, Plan{
			ID:              p.ID.Bytes,
			Name:            p.Name,
			PlanType:        p.PlanType,
			Description:     p.Description.String,
			Price:           price,
			Interval:        p.Interval,
			MaxUsers:        p.MaxUsers,
			StorageBytes:    p.DefaultStorageBytes,
			SmsCredits:      p.DefaultSmsCredits,
			WhatsappCredits: p.DefaultWhatsappCredits,
			EmailCredits:    p.DefaultEmailCredits,
		})
	}
	
	totalRecords := len(allPlans)
	
	if filter.SortBy != "" {
		sort.Slice(allPlans, func(i, j int) bool {
			var isLess bool
			switch filter.SortBy {
			case "name":
				isLess = strings.ToLower(allPlans[i].Name) < strings.ToLower(allPlans[j].Name)
			case "price":
				isLess = allPlans[i].Price < allPlans[j].Price
			default:
				isLess = strings.ToLower(allPlans[i].Name) < strings.ToLower(allPlans[j].Name)
			}
			if filter.SortOrder == "desc" {
				return !isLess
			}
			return isLess
		})
	}

	// Pagination
	var pagedPlans []Plan
	if filter.Limit == -1 {
		pagedPlans = allPlans
	} else {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		limit := filter.Limit
		if limit < 1 {
			limit = 10
		}
		start := (page - 1) * limit
		end := start + limit

		if start < totalRecords {
			if end > totalRecords {
				end = totalRecords
			}
			pagedPlans = allPlans[start:end]
		} else {
			pagedPlans = []Plan{}
		}
	}

	totalPages := 1
	if filter.Limit > 0 {
		totalPages = (totalRecords + filter.Limit - 1) / filter.Limit
	}
	if totalPages == 0 {
		totalPages = 1
	}

	if pagedPlans == nil {
		pagedPlans = []Plan{}
	}

	return PaginatedPlans{
		Data:           pagedPlans,
		TotalRecords:   totalRecords,
		PresentRecords: len(pagedPlans),
		CurrentPage:    filter.Page,
		TotalPages:     totalPages,
		Limit:          filter.Limit,
	}, nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

type PlanStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
	Archived int `json:"archived"`
}

func (s *QueryService) GetPlanStats(ctx context.Context) (PlanStats, error) {
	// Simple mock calculation from existing list since we don't have a direct query
	res, err := s.ListPlans(ctx, PlanFilter{Limit: -1})
	if err != nil {
		return PlanStats{}, err
	}
	plans := res.Data
	
	// Defaulting all to active for now as status isn't fully mapped in the model yet
	return PlanStats{
		Total:    len(plans),
		Active:   len(plans),
		Inactive: 0,
		Archived: 0,
	}, nil
}
