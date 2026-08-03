package universal

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

func (r *Repository) Create(ctx context.Context, params db.CreateUniversalParams) (db.MstUniversal, error) {
	return r.q.CreateUniversal(ctx, params)
}

func (r *Repository) GetByMasterType(ctx context.Context, masterTypeID uuid.UUID) ([]db.MstUniversal, error) {
	var pgUUID pgtype.UUID
	pgUUID.Bytes = masterTypeID
	pgUUID.Valid = true
	return r.q.GetUniversalsByMasterType(ctx, pgUUID)
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, valueName string, remarks pgtype.Text) error {
	var pgUUID pgtype.UUID
	pgUUID.Bytes = id
	pgUUID.Valid = true

	return r.q.UpdateUniversal(ctx, db.UpdateUniversalParams{
		ValueName:     valueName,
		Remarks:       remarks,
		UniversalUuid: pgUUID,
	})
}

func (r *Repository) Deactivate(ctx context.Context, id uuid.UUID) error {
	var pgUUID pgtype.UUID
	pgUUID.Bytes = id
	pgUUID.Valid = true
	return r.q.DeactivateUniversal(ctx, pgUUID)
}
