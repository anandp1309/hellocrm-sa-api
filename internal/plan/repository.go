package plan

import "hellocrm-superadmin/internal/platform/database/db"

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}
