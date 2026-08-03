package mastertype

import (
	"context"
	"time"

	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateMasterTypeRequest struct {
	Name         string `json:"name"`
	DisplayOrder int32  `json:"display_order"`
	IsSystem     bool   `json:"is_system"`
	Remarks      string `json:"remarks"`
	CreatedBy    string `json:"created_by"`
}

type UpdateMasterTypeRequest struct {
	Name         *string `json:"name"`
	DisplayOrder *int32  `json:"display_order"`
	IsSystem     *bool   `json:"is_system"`
	Remarks      *string `json:"remarks"`
	UpdatedBy    string  `json:"updated_by"`
}

func (s *Service) Create(ctx context.Context, req CreateMasterTypeRequest) (db.MstMasterType, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return db.MstMasterType{}, err
	}

	createdByUser, _ := uuid.Parse(req.CreatedBy)

	var pgID, pgCreatedBy pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	if createdByUser != uuid.Nil {
		pgCreatedBy.Bytes = createdByUser
		pgCreatedBy.Valid = true
	}

	var pgDisplayOrder pgtype.Int4
	pgDisplayOrder.Int32 = req.DisplayOrder
	pgDisplayOrder.Valid = true

	var pgIsSystem pgtype.Bool
	pgIsSystem.Bool = req.IsSystem
	pgIsSystem.Valid = true

	var pgRemarks pgtype.Text
	if req.Remarks != "" {
		pgRemarks.String = req.Remarks
		pgRemarks.Valid = true
	}

	var pgNow pgtype.Timestamptz
	pgNow.Time = time.Now()
	pgNow.Valid = true

	return s.repo.Create(ctx, db.CreateMasterTypeParams{
		MasterTypeUuid:     pgID,
		MasterTypeName:     req.Name,
		DisplayOrder:       req.DisplayOrder,
		IsSystem:           req.IsSystem,
		Remarks:            pgRemarks,
		CreatedAt:          pgNow,
		CreatedByUserUuid:  pgCreatedBy,
	})
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.MstMasterType, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int32) ([]db.MstMasterType, error) {
	if limit == 0 {
		limit = 100
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateMasterTypeRequest) (db.MstMasterType, error) {
	var pgID, pgUpdatedBy pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	updatedByUser, _ := uuid.Parse(req.UpdatedBy)
	if updatedByUser != uuid.Nil {
		pgUpdatedBy.Bytes = updatedByUser
		pgUpdatedBy.Valid = true
	}

	var pgName pgtype.Text
	if req.Name != nil {
		pgName.String = *req.Name
		pgName.Valid = true
	}

	var pgDisplayOrder pgtype.Int4
	if req.DisplayOrder != nil {
		pgDisplayOrder.Int32 = *req.DisplayOrder
		pgDisplayOrder.Valid = true
	}

	var pgIsSystem pgtype.Bool
	if req.IsSystem != nil {
		pgIsSystem.Bool = *req.IsSystem
		pgIsSystem.Valid = true
	}

	var pgRemarks pgtype.Text
	if req.Remarks != nil {
		pgRemarks.String = *req.Remarks
		pgRemarks.Valid = true
	}

	return s.repo.Update(ctx, db.UpdateMasterTypeParams{
		MasterTypeUuid:     pgID,
		MasterTypeName:     pgName, 
		DisplayOrder:       pgDisplayOrder,
		IsSystem:           pgIsSystem,
		Remarks:            pgRemarks,
		UpdatedByUserUuid:  pgUpdatedBy,
	})
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	var pgID, pgDeletedBy pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	deletedByUser, _ := uuid.Parse(deletedBy)
	if deletedByUser != uuid.Nil {
		pgDeletedBy.Bytes = deletedByUser
		pgDeletedBy.Valid = true
	}

	return s.repo.Delete(ctx, db.DeleteMasterTypeParams{
		MasterTypeUuid:     pgID,
		DeletedByUserUuid: pgDeletedBy,
	})
}
