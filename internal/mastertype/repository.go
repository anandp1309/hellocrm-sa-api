package mastertype

import (
	"context"

	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

func (r *Repository) Create(ctx context.Context, params db.CreateMasterTypeParams) (db.MstMasterType, error) {
	return r.q.CreateMasterType(ctx, params)
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (db.MstMasterType, error) {
	var pgUUID pgtype.UUID
	pgUUID.Bytes = id
	pgUUID.Valid = true
	return r.q.GetMasterType(ctx, pgUUID)
}

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]db.MstMasterType, error) {
	return r.q.ListMasterTypes(ctx, db.ListMasterTypesParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *Repository) Update(ctx context.Context, params db.UpdateMasterTypeParams) (db.MstMasterType, error) {
	return r.q.UpdateMasterType(ctx, params)
}

func (r *Repository) Delete(ctx context.Context, params db.DeleteMasterTypeParams) error {
	return r.q.DeleteMasterType(ctx, params)
}
