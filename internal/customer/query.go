package customer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

type Customer struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (s *QueryService) GetCustomerByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	if id == uuid.Nil {
		return Customer{}, fmt.Errorf("invalid id")
	}
	// STUB: Replace with actual database call
	return Customer{
		ID:    id,
		Name:  "Stub Customer",
		Email: "stub@example.com",
	}, nil
}

func (s *QueryService) ListCustomers(ctx context.Context) ([]Customer, error) {
	// STUB: Replace with actual database call
	return []Customer{
		{
			ID:    uuid.Must(uuid.NewV7()),
			Name:  "Stub Customer 1",
			Email: "stub1@example.com",
		},
	}, nil
}
