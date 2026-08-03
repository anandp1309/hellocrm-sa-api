package tenant

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CommandService owns business rules and write operations for tenants.
type CommandService struct {
	repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{
		repo: repo,
	}
}

type CreateTenantParams struct {
	Name              string
	Email             string
	ContactPersonName string
	MobileNumber      string
	CountryName       string
	StateName         string
	CityName          string
	Address           string
	GstNumber         string
	Remarks           string
}

func (s *CommandService) CreateTenant(ctx context.Context, params CreateTenantParams) (uuid.UUID, error) {
	// 1. Validate domain invariants
	if params.Name == "" {
		return uuid.Nil, fmt.Errorf("tenant name is required")
	}

	// 2. Generate UUIDv7 as per rules
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate UUIDv7: %w", err)
	}

	// 3. Call repository
	var pgID pgtype.UUID
	err = pgID.Scan(id.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID: %w", err)
	}

	err = s.repo.CreateTenant(ctx, pgID, params)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "uk_tenant_code") {
			return uuid.Nil, fmt.Errorf("Customer already exists")
		}
		return uuid.Nil, fmt.Errorf("failed to insert tenant: %w", err)
	}

	return id, nil
}

type UpdateTenantParams struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func (s *CommandService) UpdateTenant(ctx context.Context, params UpdateTenantParams) error {
	if params.Name == "" {
		return fmt.Errorf("tenant name is required")
	}

	var pgID pgtype.UUID
	err := pgID.Scan(params.ID.String())
	if err != nil {
		return fmt.Errorf("failed to parse UUID: %w", err)
	}

	return s.repo.UpdateTenant(ctx, pgID, params.Name, params.Email)
}

func (s *CommandService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	var pgID pgtype.UUID
	err := pgID.Scan(id.String())
	if err != nil {
		return fmt.Errorf("failed to parse UUID: %w", err)
	}

	return s.repo.DeleteTenant(ctx, pgID)
}
