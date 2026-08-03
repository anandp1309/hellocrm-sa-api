package gateway

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

func (r *Repository) ListPaginated(ctx context.Context, search, gateway, status, startDate, endDate string, limit, offset int32) ([]db.ListGatewayTransactionsPaginatedRow, error) {
	return r.q.ListGatewayTransactionsPaginated(ctx, db.ListGatewayTransactionsPaginatedParams{
		Column1: search,
		Column2: gateway,
		Column3: status,
		Column6: startDate,
		Column7: endDate,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) GetStats(ctx context.Context) (db.GetGatewayStatsRow, error) {
	return r.q.GetGatewayStats(ctx)
}
