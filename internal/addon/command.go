package addon

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CommandService struct {
	repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{repo: repo}
}

type CreateAddonParams struct {
	Name       string
	Category   string
	Limit      string
	Price      string
	Status     string
	IconSvg    string
	IconColor  string
}

func (s *CommandService) CreateAddon(ctx context.Context, params CreateAddonParams) (uuid.UUID, error) {
	id := uuid.New()
	_, err := s.repo.q.CreateAddon(ctx, db.CreateAddonParams{
		AddonUuid:  pgtype.UUID{Bytes: id, Valid: true},
		Name:       params.Name,
		Category:   params.Category,
		AddonLimit: params.Limit,
		Price:      params.Price,
		Status:     params.Status,
		IconSvg:    pgtype.Text{String: params.IconSvg, Valid: true},
		IconColor:  pgtype.Text{String: params.IconColor, Valid: true},
	})
	return id, err
}

func (s *CommandService) DeleteAddon(ctx context.Context, id uuid.UUID) error {
	status := "Inactive"
	return s.UpdateAddon(ctx, UpdateAddonParams{
		ID:     id,
		Status: &status,
	})
}

type UpdateAddonParams struct {
	ID        uuid.UUID
	Name      *string
	Category  *string
	Limit     *string
	Price     *string
	Status    *string
	IconSvg   *string
	IconColor *string
}

func (s *CommandService) UpdateAddon(ctx context.Context, params UpdateAddonParams) error {
	arg := db.UpdateAddonParams{
		AddonUuid: pgtype.UUID{Bytes: params.ID, Valid: true},
	}
	if params.Name != nil {
		arg.Name = pgtype.Text{String: *params.Name, Valid: true}
	}
	if params.Category != nil {
		arg.Category = pgtype.Text{String: *params.Category, Valid: true}
	}
	if params.Limit != nil {
		arg.AddonLimit = pgtype.Text{String: *params.Limit, Valid: true}
	}
	if params.Price != nil {
		arg.Price = pgtype.Text{String: *params.Price, Valid: true}
	}
	if params.Status != nil {
		arg.Status = pgtype.Text{String: *params.Status, Valid: true}
	}
	if params.IconSvg != nil {
		arg.IconSvg = pgtype.Text{String: *params.IconSvg, Valid: true}
	}
	if params.IconColor != nil {
		arg.IconColor = pgtype.Text{String: *params.IconColor, Valid: true}
	}

	_, err := s.repo.q.UpdateAddon(ctx, arg)
	return err
}
