package customer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CommandService struct {
	repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{repo: repo}
}

type CreateCustomerParams struct {
	Name  string
	Email string
}

type UpdateCustomerParams struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func (s *CommandService) CreateCustomer(ctx context.Context, params CreateCustomerParams) (uuid.UUID, error) {
	if params.Name == "" {
		return uuid.Nil, fmt.Errorf("name is required")
	}
	id, _ := uuid.NewV7()
	// STUB: Replace with actual database call
	return id, nil
}

func (s *CommandService) UpdateCustomer(ctx context.Context, params UpdateCustomerParams) error {
	if params.ID == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	// STUB: Replace with actual database call
	return nil
}

func (s *CommandService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	// STUB: Replace with actual database call
	return nil
}
