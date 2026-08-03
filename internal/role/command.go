package role

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CommandService struct {
	repo Repository
}

func NewCommandService(repo Repository) *CommandService {
	return &CommandService{repo: repo}
}

type UpdateRoleRequest struct {
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
}

func (s *CommandService) UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleRequest) error {
	var pgID pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	var remarks pgtype.Text
	remarks.String = req.Remarks
	remarks.Valid = req.Remarks != ""

	return s.repo.UpdateRole(ctx, db.UpdateRoleParams{
		RoleUuid: pgID,
		RoleName: req.Name,
		Remarks:  remarks,
	})
}
