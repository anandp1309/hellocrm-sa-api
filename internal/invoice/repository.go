package invoice

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

func (r *Repository) ListPaginated(ctx context.Context, search, status, startDate, endDate string, limit, offset int32) ([]db.ListInvoicesPaginatedRow, error) {
	return r.q.ListInvoicesPaginated(ctx, db.ListInvoicesPaginatedParams{
		Column1: search,
		Column2: status,
		Column5: startDate,
		Column6: endDate,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) GetStats(ctx context.Context) (db.GetInvoiceStatsRow, error) {
	return r.q.GetInvoiceStats(ctx)
}
