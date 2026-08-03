package tenant

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"hellocrm-superadmin/internal/platform/database/db"
)

// Repository encapsulates data access for the tenant module.
// Following the OS rule: avoid generic repositories, use concrete ones.
type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) CreateTenant(ctx context.Context, id pgtype.UUID, params CreateTenantParams) error {
		var emailPg pgtype.Text
	if params.Email != "" {
		emailPg = pgtype.Text{String: params.Email, Valid: true}
	}
	idStr := fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
	
	// Convert optional string fields to pgtype.Text
	toPgText := func(s string) pgtype.Text {
		if s != "" {
			return pgtype.Text{String: s, Valid: true}
		}
		return pgtype.Text{}
	}

	return r.queries.CreateTenant(ctx, db.CreateTenantParams{
		TenantUuid:         id,
		TenantCode:         params.Name[:min(3, len(params.Name))],
		TenantID:           idStr,
		TenantName:         params.Name,
		EmailAddress:       emailPg,
		ContactPersonName:  toPgText(params.ContactPersonName),
		MobileNumber:       toPgText(params.MobileNumber),
		CountryName:        toPgText(params.CountryName),
		StateName:          toPgText(params.StateName),
		CityName:           toPgText(params.CityName),
		Address:            toPgText(params.Address),
		GstNumber:          toPgText(params.GstNumber),
		Remarks:            toPgText(params.Remarks),
	})
}

func (r *Repository) GetTenantByID(ctx context.Context, id pgtype.UUID) (db.GetTenantByIDRow, error) {
	return r.queries.GetTenantByID(ctx, id)
}

func (r *Repository) ListTenantsPaginated(ctx context.Context, search, status, planType, plan, billingCycle string, limit, offset int32) ([]db.ListTenantsPaginatedRow, error) {
	return r.queries.ListTenantsPaginated(ctx, db.ListTenantsPaginatedParams{
		Column1: search,
		Column2: status,
		Column3: planType,
		Column4: plan,
		Column5: billingCycle,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) GetTenantStats(ctx context.Context) (db.GetTenantStatsRow, error) {
	return r.queries.GetTenantStats(ctx)
}

func (r *Repository) UpdateTenant(ctx context.Context, id pgtype.UUID, name, email string) error {
	var emailPg pgtype.Text
	if email != "" {
		emailPg = pgtype.Text{String: email, Valid: true}
	}
	return r.queries.UpdateTenant(ctx, db.UpdateTenantParams{
		TenantUuid:   id,
		TenantName:   name,
		EmailAddress: emailPg,
	})
}

func (r *Repository) DeleteTenant(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteTenant(ctx, id)
}
