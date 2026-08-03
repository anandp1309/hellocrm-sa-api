package usage

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"
)

type QueryService interface {
	GetUsageStats(ctx context.Context) (db.GetUsageStatsRow, error)
}

type queryService struct {
	queries *db.Queries
}

func NewQueryService(queries *db.Queries) QueryService {
	return &queryService{queries: queries}
}

func (s *queryService) GetUsageStats(ctx context.Context) (db.GetUsageStatsRow, error) {
	return s.queries.GetUsageStats(ctx)
}
