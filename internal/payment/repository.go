package payment

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{
		q: q,
	}
}

func (r *Repository) ListPaginated(ctx context.Context, search, status, method, plan, startDate, endDate string, limit, offset int32) ([]db.ListPaymentsPaginatedRow, error) {
	return r.q.ListPaymentsPaginated(ctx, db.ListPaymentsPaginatedParams{
		Column1: search,
		Column2: status,
		Column3: method,
		Column4: plan,
		Column7: startDate,
		Column8: endDate,
		Offset:  offset,
		Limit:   limit,
	})
}

func (r *Repository) GetStats(ctx context.Context) (db.GetPaymentStatsRow, error) {
	return r.q.GetPaymentStats(ctx)
}

func (r *Repository) Create(ctx context.Context, arg db.CreatePaymentParams) (db.TenantSubscriptionPayment, error) {
	return r.q.CreatePayment(ctx, arg)
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (db.TenantSubscriptionPayment, error) {
	var pgUUID pgtype.UUID
	pgUUID.Bytes = id
	pgUUID.Valid = true
	return r.q.GetPayment(ctx, pgUUID)
}

func (r *Repository) Update(ctx context.Context, arg db.UpdatePaymentParams) (db.TenantSubscriptionPayment, error) {
	return r.q.UpdatePayment(ctx, arg)
}

func (r *Repository) Delete(ctx context.Context, arg db.DeletePaymentParams) error {
	return r.q.DeletePayment(ctx, arg)
}
