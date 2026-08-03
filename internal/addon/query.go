package addon

import (
	"context"
	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) ListAddons(ctx context.Context) ([]db.MstAddon, error) {
	return s.repo.q.ListAddons(ctx)
}

func (s *QueryService) GetAddonByID(ctx context.Context, id uuid.UUID) (db.MstAddon, error) {
	return s.repo.q.GetAddon(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
