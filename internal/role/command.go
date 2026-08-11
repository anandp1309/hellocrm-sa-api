package role

import (
	"context"
	"fmt"
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

type PermissionPayload struct {
	ModuleUuid string `json:"moduleUuid"`
	Module     string `json:"module"`
	RoleUuid   string `json:"roleUuid"`
	TenantUuid string `json:"tenantUuid"`
	V          bool   `json:"v"`
	V_Uuid     string `json:"v_uuid"`
	C          bool   `json:"c"`
	C_Uuid     string `json:"c_uuid"`
	E          bool   `json:"e"`
	E_Uuid     string `json:"e_uuid"`
	D          bool   `json:"d"`
	D_Uuid     string `json:"d_uuid"`
	X          bool   `json:"x"`
	X_Uuid     string `json:"x_uuid"`
	I          bool   `json:"i"`
	I_Uuid     string `json:"i_uuid"`
	M          bool   `json:"m"`
	M_Uuid     string `json:"m_uuid"`
}

type CreateRoleRequest struct {
	Name        string              `json:"name"`
	Remarks     string              `json:"remarks"`
	Permissions []PermissionPayload `json:"permissions"`
}

type UpdateRolePermissionsRequest struct {
	RoleId      string              `json:"roleId"`
	Permissions []PermissionPayload `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
}

func (s *CommandService) CreateRole(ctx context.Context, tenantID uuid.UUID, req CreateRoleRequest) error {
	roleID := uuid.New()
	var pgRoleID pgtype.UUID
	pgRoleID.Bytes = roleID
	pgRoleID.Valid = true

	var pgTenantID pgtype.UUID
	pgTenantID.Bytes = tenantID
	pgTenantID.Valid = true

	var remarks pgtype.Text
	remarks.String = req.Remarks
	remarks.Valid = req.Remarks != ""

	_, err := s.repo.CreateNewRole(ctx, db.CreateNewRoleParams{
		RoleUuid:   pgRoleID,
		TenantUuid: pgTenantID,
		RoleName:   req.Name,
		Remarks:    remarks,
	})
	if err != nil {
		return err
	}

	return s.savePermissions(ctx, pgRoleID, req.Permissions)
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

func (s *CommandService) UpdateRolePermissions(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateRolePermissionsRequest) error {
	var pgID pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	return s.savePermissions(ctx, pgID, req.Permissions)
}

func (s *CommandService) savePermissions(ctx context.Context, roleUuid pgtype.UUID, permissions []PermissionPayload) error {
	for _, p := range permissions {
		roleStr := fmt.Sprintf("%x-%x-%x-%x-%x", roleUuid.Bytes[0:4], roleUuid.Bytes[4:6], roleUuid.Bytes[6:8], roleUuid.Bytes[8:10], roleUuid.Bytes[10:16])
		fmt.Printf("Processing module %s: v=%v (uuid=%s) for role=%s\n", p.Module, p.V, p.V_Uuid, roleStr)
		actions := []struct {
			Value bool
			Uuid  string
		}{
			{p.V, p.V_Uuid},
			{p.C, p.C_Uuid},
			{p.E, p.E_Uuid},
			{p.D, p.D_Uuid},
			{p.X, p.X_Uuid},
			{p.I, p.I_Uuid},
			{p.M, p.M_Uuid},
		}

		for _, a := range actions {
			if a.Uuid != "" {
				var permUuid pgtype.UUID
				err := permUuid.Scan(a.Uuid)
				if err != nil {
					continue
				}
				
				if a.Value {
					var rpUuid pgtype.UUID
					rpUuid.Bytes = uuid.New()
					rpUuid.Valid = true
					
					var modUuid pgtype.UUID
					err := modUuid.Scan(p.ModuleUuid)
					if err != nil {
						continue
					}

					err = s.repo.AddRolePermission(ctx, db.AddRolePermissionParams{
						RolePermissionUuid: rpUuid,
						RoleUuid:           roleUuid,
						ModuleUuid:         modUuid,
						PermissionUuid:     permUuid,
					})
					if err != nil {
						fmt.Println("Error adding permission:", err)
					}
				} else {
					var modUuid pgtype.UUID
					err := modUuid.Scan(p.ModuleUuid)
					if err != nil {
						continue
					}

					err = s.repo.DeleteRolePermission(ctx, db.DeleteRolePermissionParams{
						RoleUuid:       roleUuid,
						ModuleUuid:     modUuid,
						PermissionUuid: permUuid,
					})
					if err != nil {
						fmt.Println("Error deleting permission:", err)
					}
				}
			}
		}
	}
	return nil
}
