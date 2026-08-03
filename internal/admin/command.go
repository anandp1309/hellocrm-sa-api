package admin

import (
	"context"
	"fmt"
	"hellocrm-superadmin/internal/platform/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{repo: repo}
}

type CreateAdminRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RoleID    string `json:"role_id"`
	StatusID  string `json:"status_id"`
}

func (s *CommandService) CreateAdmin(ctx context.Context, req CreateAdminRequest) (uuid.UUID, error) {
	// 1. Get Tenant UUID
	var tenantUUID pgtype.UUID
	
	// Because admin is super admin, we need the tenant SYS-001
	tenantRow, err := s.repo.q.GetTenantByCode(ctx, "SYS-001")
	if err == nil {
		tenantUUID = tenantRow
	} else {
		// fallback to searching directly, assuming tenant exists and code was wrong
		// but since SYS-001 always exists, this should not be hit.
		return uuid.Nil, fmt.Errorf("tenant SYS-001 not found")
	}

	userId, _ := uuid.NewV7()
	var pgUserId pgtype.UUID
	pgUserId.Bytes = userId
	pgUserId.Valid = true

	authId, _ := uuid.NewV7()
	var pgAuthId pgtype.UUID; pgAuthId.Bytes = authId; pgAuthId.Valid = true

	utId, _ := uuid.NewV7()
	var pgUtId pgtype.UUID; pgUtId.Bytes = utId; pgUtId.Valid = true

	uwId, _ := uuid.NewV7()
	var pgUwId pgtype.UUID; pgUwId.Bytes = uwId; pgUwId.Valid = true

	var roleUUID pgtype.UUID; rId, _ := uuid.Parse(req.RoleID); roleUUID.Bytes = rId; roleUUID.Valid = true
	var statusUUID pgtype.UUID; stId, _ := uuid.Parse(req.StatusID); statusUUID.Bytes = stId; statusUUID.Valid = true

	var pgLastName pgtype.Text; if req.LastName != "" { pgLastName.String = req.LastName; pgLastName.Valid = true }
	var pgEmail pgtype.Text; if req.Email != "" { pgEmail.String = req.Email; pgEmail.Valid = true }
	var nullUUID pgtype.UUID

	err = s.repo.q.CreateAdminUser(ctx, db.CreateAdminUserParams{
		UserUuid: pgUserId, FirstName: req.FirstName, LastName: pgLastName, EmailAddress: pgEmail, UserStatusUniversalUuid: statusUUID, CreatedByUserUuid: nullUUID,
	})
	if err != nil { return uuid.Nil, fmt.Errorf("user: %w", err) }

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	err = s.repo.q.CreateAdminAuth(ctx, db.CreateAdminAuthParams{
		UserAuthUuid: pgAuthId, TenantUuid: tenantUUID, UserUuid: pgUserId, LoginID: req.Username, PasswordHash: string(hashed),
	})
	if err != nil { return uuid.Nil, fmt.Errorf("auth: %w", err) }

	var empCode pgtype.Text
	empCode.String = "SA-NEW"
	empCode.Valid = true
	err = s.repo.q.CreateAdminTenant(ctx, db.CreateAdminTenantParams{
		UserTenantUuid: pgUtId, TenantUuid: tenantUUID, UserUuid: pgUserId, EmployeeCode: empCode,
	})
	if err != nil { return uuid.Nil, fmt.Errorf("tenant: %w", err) }

	var wsUUID pgtype.UUID
	err = s.repo.q.CreateAdminWorkspace(ctx, db.CreateAdminWorkspaceParams{
		UserWorkspaceUuid: pgUwId, TenantUuid: tenantUUID, UserTenantUuid: pgUtId, WorkspaceUuid: wsUUID, RoleUuid: roleUUID,
	})
	if err != nil { return uuid.Nil, fmt.Errorf("workspace: %w", err) }

	return userId, nil
}

type UpdateAdminRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RoleID    string `json:"role_id"`
	StatusID  string `json:"status_id"`
	MfaEnabled bool   `json:"mfa_enabled"`
}

func (s *CommandService) UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error {
	var pgUserId pgtype.UUID
	pgUserId.Bytes = id
	pgUserId.Valid = true

	var roleUUID pgtype.UUID; rId, _ := uuid.Parse(req.RoleID); roleUUID.Bytes = rId; roleUUID.Valid = true
	var statusUUID pgtype.UUID; stId, _ := uuid.Parse(req.StatusID); statusUUID.Bytes = stId; statusUUID.Valid = true
	var pgLastName pgtype.Text; if req.LastName != "" { pgLastName.String = req.LastName; pgLastName.Valid = true }
	var pgEmail pgtype.Text; if req.Email != "" { pgEmail.String = req.Email; pgEmail.Valid = true }
	var nullUUID pgtype.UUID

	err := s.repo.q.UpdateAdminUser(ctx, db.UpdateAdminUserParams{
		UserUuid: pgUserId, FirstName: req.FirstName, LastName: pgLastName, EmailAddress: pgEmail, UserStatusUniversalUuid: statusUUID, UpdatedByUserUuid: nullUUID,
	})
	if err != nil { return err }

	var hash string
	if req.Password != "" {
		h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hash = string(h)
	}

	err = s.repo.q.UpdateAdminAuth(ctx, db.UpdateAdminAuthParams{
		UserUuid: pgUserId, LoginID: req.Username, Column3: hash, IsMobileVerified: req.MfaEnabled,
	})
	if err != nil { return err }

	err = s.repo.q.UpdateAdminWorkspace(ctx, db.UpdateAdminWorkspaceParams{
		UserUuid: pgUserId, RoleUuid: roleUUID,
	})
	return err
}

func (s *CommandService) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	var pgUserId pgtype.UUID
	pgUserId.Bytes = id
	pgUserId.Valid = true
	var nullUUID pgtype.UUID
	return s.repo.q.DeleteAdmin(ctx, db.DeleteAdminParams{UserUuid: pgUserId, UpdatedByUserUuid: nullUUID})
}
