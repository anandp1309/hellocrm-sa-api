package identity

import (
	"hellocrm-superadmin/internal/platform/database/db"
)

// Repository encapsulates data access for identity (users, roles).
type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}
