package role

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	ListRoles(ctx context.Context) ([]db.ListRolesRow, error)
	GetRolePermissions(ctx context.Context, roleUuid pgtype.UUID) ([]db.GetRolePermissionsRow, error)
	UpdateRole(ctx context.Context, arg db.UpdateRoleParams) error
}

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) Repository {
	return &repository{db: db}
}

func (r *repository) ListRoles(ctx context.Context) ([]db.ListRolesRow, error) {
	return r.db.ListRoles(ctx)
}

func (r *repository) GetRolePermissions(ctx context.Context, roleUuid pgtype.UUID) ([]db.GetRolePermissionsRow, error) {
	return r.db.GetRolePermissions(ctx, roleUuid)
}

func (r *repository) UpdateRole(ctx context.Context, arg db.UpdateRoleParams) error {
	return r.db.UpdateRole(ctx, arg)
}
