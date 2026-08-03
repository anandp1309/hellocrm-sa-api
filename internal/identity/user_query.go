package identity

import (
	"context"

	"github.com/google/uuid"
	"hellocrm-superadmin/internal/auth"
)

// UserQueryService handles reading user data.
type UserQueryService struct {
	repo *Repository
}

func NewUserQueryService(repo *Repository) *UserQueryService {
	return &UserQueryService{
		repo: repo,
	}
}

type UserView struct {
	ID       uuid.UUID  `json:"id"`
	Email    string     `json:"email"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	RoleID   uuid.UUID  `json:"role_id"`
}

func (s *UserQueryService) GetUserByID(ctx context.Context, actor auth.Actor, id uuid.UUID) (UserView, error) {
	// Authorization Check
	if !actor.HasPermission("users:read") && actor.UserID != id {
		return UserView{}, auth.ErrUnauthorized
	}

	// Fetch data via repo (SQLC mock)
	// user, err := s.repo.queries.GetUserByID(ctx, id)

	return UserView{
		ID:    id,
		Email: "example@hellocrm.com",
	}, nil
}
