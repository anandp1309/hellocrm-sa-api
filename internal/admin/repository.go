package admin

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

func (r *Repository) ListPaginated(ctx context.Context, search, role, status, mfa string, limit, offset int32) ([]db.ListAdminsPaginatedRow, error) {
	return r.q.ListAdminsPaginated(ctx, db.ListAdminsPaginatedParams{
		Column1: search,
		Column2: role,
		Column3: status,
		Column4: mfa,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) GetStats(ctx context.Context) (db.GetAdminStatsRow, error) {
	return r.q.GetAdminStats(ctx)
}
