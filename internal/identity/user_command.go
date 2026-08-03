package identity

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"hellocrm-superadmin/internal/auth"
)

// UserCommandService handles business logic and writes for users.
type UserCommandService struct {
	repo *Repository
}

func NewUserCommandService(repo *Repository) *UserCommandService {
	return &UserCommandService{
		repo: repo,
	}
}

type CreateUserParams struct {
	Email    string
	Password string
	RoleID   uuid.UUID
	TenantID *uuid.UUID // Can be nil for superadmins
}

// CreateUser executes the command to create a new user. Requires actor context.
func (s *UserCommandService) CreateUser(ctx context.Context, actor auth.Actor, params CreateUserParams) (uuid.UUID, error) {
	// Authorization Check: Does the actor have permission to create users?
	if !actor.HasPermission("users:write") {
		return uuid.Nil, auth.ErrUnauthorized
	}

	// 1. Validate domain invariants
	if params.Email == "" {
		return uuid.Nil, fmt.Errorf("email is required")
	}
	if len(params.Password) < 8 {
		return uuid.Nil, fmt.Errorf("password must be at least 8 characters")
	}
	if params.RoleID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("role ID is required")
	}

	// 2. Generate UUIDv7
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate UUIDv7: %w", err)
	}

	// 3. Hash password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}
	hashedPassword := string(hashedBytes)

	// 4. Save to DB (using repository wrapper when sqlc is generated)
	// err = s.repo.queries.CreateUser(ctx, db.CreateUserParams{
	// 	ID:           id,
	// 	Email:        params.Email,
	// 	PasswordHash: hashedPassword,
	// 	TenantID:     params.TenantID,
	// 	RoleID:       params.RoleID,
	// })
	
	_ = hashedPassword // Keep compiler happy

	return id, nil
}
