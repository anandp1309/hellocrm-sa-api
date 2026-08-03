package addon

import (
	"hellocrm-superadmin/internal/platform/database/db"
)

type Repository struct {
	q  *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		q:  queries,
	}
}
