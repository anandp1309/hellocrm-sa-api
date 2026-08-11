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
	ListModules(ctx context.Context) ([]db.ListModulesRow, error)
	ListAllPermissions(ctx context.Context) ([]db.ListAllPermissionsRow, error)
	CreateNewRole(ctx context.Context, arg db.CreateNewRoleParams) (pgtype.UUID, error)
	GetPermissionByCode(ctx context.Context, permissionCode string) (pgtype.UUID, error)
	AddRolePermission(ctx context.Context, arg db.AddRolePermissionParams) error
	DeleteRolePermission(ctx context.Context, arg db.DeleteRolePermissionParams) error
	DeleteRolePermissions(ctx context.Context, roleUuid pgtype.UUID) error
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

func (r *repository) ListModules(ctx context.Context) ([]db.ListModulesRow, error) {
	return r.db.ListModules(ctx)
}

func (r *repository) ListAllPermissions(ctx context.Context) ([]db.ListAllPermissionsRow, error) {
	return r.db.ListAllPermissions(ctx)
}

func (r *repository) CreateNewRole(ctx context.Context, arg db.CreateNewRoleParams) (pgtype.UUID, error) {
	return r.db.CreateNewRole(ctx, arg)
}

func (r *repository) GetPermissionByCode(ctx context.Context, permissionCode string) (pgtype.UUID, error) {
	return r.db.GetPermissionByCode(ctx, permissionCode)
}

func (r *repository) AddRolePermission(ctx context.Context, arg db.AddRolePermissionParams) error {
	return r.db.AddRolePermission(ctx, arg)
}

func (r *repository) DeleteRolePermission(ctx context.Context, arg db.DeleteRolePermissionParams) error {
	return r.db.DeleteRolePermission(ctx, arg)
}

func (r *repository) DeleteRolePermissions(ctx context.Context, roleUuid pgtype.UUID) error {
	return r.db.DeleteRolePermissions(ctx, roleUuid)
}
