package dashboard

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"
)

type QueryService interface {
	GetDashboardStats(ctx context.Context) (db.GetDashboardStatsRow, error)
}

type queryService struct {
	queries *db.Queries
}

func NewQueryService(queries *db.Queries) QueryService {
	return &queryService{queries: queries}
}

func (s *queryService) GetDashboardStats(ctx context.Context) (db.GetDashboardStatsRow, error) {
	return s.queries.GetDashboardStats(ctx)
}
