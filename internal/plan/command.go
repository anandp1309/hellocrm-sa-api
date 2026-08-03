package plan

import (
	"context"
	"fmt"
	"strconv"

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

type CreatePlanParams struct {
	Name            string
	Description     string
	Price           float64
	Interval        string
	MaxUsers        int32
	StorageBytes    int64
	SmsCredits      int32
	WhatsappCredits int32
	EmailCredits    int32
}

type UpdatePlanParams struct {
	ID              uuid.UUID
	Name            string
	Description     string
	Price           float64
	Interval        string
	MaxUsers        int32
	StorageBytes    int64
	SmsCredits      int32
	WhatsappCredits int32
	EmailCredits    int32
}

func toPgNumeric(val float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(strconv.FormatFloat(val, 'f', -1, 64))
	return num
}

func (s *CommandService) CreatePlan(ctx context.Context, params CreatePlanParams) (uuid.UUID, error) {
	if params.Name == "" {
		return uuid.Nil, fmt.Errorf("name is required")
	}
	id, _ := uuid.NewV7()
	
	maxUsers := params.MaxUsers
	if maxUsers <= 0 {
		maxUsers = -1
	}

	err := s.repo.queries.CreateMstPlan(ctx, db.CreateMstPlanParams{
		PlanUuid:               pgtype.UUID{Bytes: id, Valid: true},
		PlanName:               params.Name,
		Remarks:                pgtype.Text{String: params.Description, Valid: params.Description != ""},
		MaxUsers:               maxUsers,
		DefaultStorageBytes:    params.StorageBytes,
		DefaultSmsCredits:      params.SmsCredits,
		DefaultWhatsappCredits: params.WhatsappCredits,
		DefaultEmailCredits:    params.EmailCredits,
	})
	if err != nil {
		return uuid.Nil, err
	}

	priceId, _ := uuid.NewV7()
	err = s.repo.queries.CreateMstPlanPrice(ctx, db.CreateMstPlanPriceParams{
		PlanPriceUuid: pgtype.UUID{Bytes: priceId, Valid: true},
		PlanUuid:      pgtype.UUID{Bytes: id, Valid: true},
		PriceAmount:   toPgNumeric(params.Price),
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *CommandService) UpdatePlan(ctx context.Context, params UpdatePlanParams) error {
	if params.ID == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	
	maxUsers := params.MaxUsers
	if maxUsers <= 0 {
		maxUsers = -1
	}

	err := s.repo.queries.UpdateMstPlan(ctx, db.UpdateMstPlanParams{
		PlanUuid:               pgtype.UUID{Bytes: params.ID, Valid: true},
		PlanName:               params.Name,
		Remarks:                pgtype.Text{String: params.Description, Valid: params.Description != ""},
		MaxUsers:               maxUsers,
		DefaultStorageBytes:    params.StorageBytes,
		DefaultSmsCredits:      params.SmsCredits,
		DefaultWhatsappCredits: params.WhatsappCredits,
		DefaultEmailCredits:    params.EmailCredits,
	})
	if err != nil {
		return err
	}

	err = s.repo.queries.UpdateMstPlanPrice(ctx, db.UpdateMstPlanPriceParams{
		PlanUuid:    pgtype.UUID{Bytes: params.ID, Valid: true},
		PriceAmount: toPgNumeric(params.Price),
	})
	return err
}

func (s *CommandService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	
	return s.repo.queries.DeleteMstPlan(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
