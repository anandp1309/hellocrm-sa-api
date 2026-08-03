package payment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"hellocrm-superadmin/internal/platform/database/db"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{
		repo: repo,
	}
}

type PaginatedPayments struct {
	Data       []PaymentListView `json:"data"`
	Total      int64             `json:"total"`
	Page       int32             `json:"page"`
	Limit      int32             `json:"limit"`
	TotalPages int64             `json:"totalPages"`
}

type PaymentListView struct {
	ID            string  `json:"id"`
	TransactionId string  `json:"transactionId"`
	Customer      string  `json:"customer"`
	Email         string  `json:"email"`
	Status        string  `json:"status"`
	Plan          string  `json:"plan"`
	Amount        string  `json:"amount"`
	Method        string  `json:"method"`
	MethodIcon    string  `json:"methodIcon"`
	Date          string  `json:"date"`
	Invoice       string  `json:"invoice"`
	Initials      string  `json:"initials"`
	IconBg        string  `json:"iconBg"`
	IconColor     string  `json:"iconColor"`
}

func (s *QueryService) ListPaymentsPaginated(ctx context.Context, search, status, method, plan, startDate, endDate string, page, limit int32) (PaginatedPayments, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := s.repo.ListPaginated(ctx, search, status, method, plan, startDate, endDate, limit, offset)
	if err != nil {
		return PaginatedPayments{}, err
	}

	var results []PaymentListView
	var total int64 = 0

	for _, row := range rows {
		total = row.TotalRecords
		rowID, _ := uuid.Parse(fmt.Sprintf("%x-%x-%x-%x-%x", row.TenantSubscriptionPaymentUuid.Bytes[0:4], row.TenantSubscriptionPaymentUuid.Bytes[4:6], row.TenantSubscriptionPaymentUuid.Bytes[6:8], row.TenantSubscriptionPaymentUuid.Bytes[8:10], row.TenantSubscriptionPaymentUuid.Bytes[10:16]))

		amountStr := "₹ 0"
		if row.Amount.Valid {
			f, _ := row.Amount.Float64Value()
			amountStr = fmt.Sprintf("₹ %.0f", f.Float64)
		}

		dateStr := ""
		if row.PaymentDate.Valid {
			dateStr = row.PaymentDate.Time.Format("02 Jan 2006 03:04 PM")
		}

		// Calculate UI cosmetics
		initials := ""
		if len(row.Customer) >= 2 {
			initials = strings.ToUpper(row.Customer[:2])
		} else if len(row.Customer) == 1 {
			initials = strings.ToUpper(row.Customer[:1])
		}

		var iconBg, iconColor string
		colors := []struct{ bg, fg string }{
			{"bg-blue-50", "text-blue-600"},
			{"bg-green-50", "text-green-600"},
			{"bg-purple-50", "text-purple-600"},
			{"bg-orange-50", "text-orange-600"},
			{"bg-pink-50", "text-pink-600"},
		}
		if len(row.Customer) > 0 {
			idx := int(row.Customer[0]) % len(colors)
			iconBg = colors[idx].bg
			iconColor = colors[idx].fg
		} else {
			iconBg = "bg-gray-50"
			iconColor = "text-gray-600"
		}

		methodIcon := "none"
		if strings.Contains(strings.ToLower(row.Method), "visa") {
			methodIcon = "visa"
		} else if strings.Contains(strings.ToLower(row.Method), "mastercard") {
			methodIcon = "mastercard"
		} else if strings.Contains(strings.ToLower(row.Method), "upi") {
			methodIcon = "upi"
		} else if strings.Contains(strings.ToLower(row.Method), "net banking") {
			methodIcon = "bank"
		}

		results = append(results, PaymentListView{
			ID:            rowID.String(),
			TransactionId: row.TransactionID,
			Customer:      row.Customer,
			Email:         row.EmailAddress.String,
			Status:        row.Status,
			Plan:          row.Plan,
			Amount:        amountStr,
			Method:        row.Method,
			MethodIcon:    methodIcon,
			Date:          dateStr,
			Invoice:       "INV-" + row.TransactionID, // mock invoice based on txn ID
			Initials:      initials,
			IconBg:        iconBg,
			IconColor:     iconColor,
		})
	}

	totalPages := int64(1)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return PaginatedPayments{
		Data:       results,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

type PaymentStats struct {
	TotalPayments  int64  `json:"totalPayments"`
	TotalReceived  string `json:"totalReceived"`
	PendingAmount  string `json:"pendingAmount"`
	PendingCount   int64  `json:"pendingCount"`
	FailedAmount   string `json:"failedAmount"`
	FailedCount    int64  `json:"failedCount"`
	RefundedAmount string `json:"refundedAmount"`
	RefundedCount  int64  `json:"refundedCount"`
}

func (s *QueryService) GetPaymentStats(ctx context.Context) (PaymentStats, error) {
	row, err := s.repo.GetStats(ctx)
	if err != nil {
		return PaymentStats{}, err
	}

	formatAmt := func(v interface{}) string {
		if v == nil {
			return "₹ 0"
		}
		switch val := v.(type) {
		case float64:
			return fmt.Sprintf("₹ %.0f", val)
		case float32:
			return fmt.Sprintf("₹ %.0f", val)
		case int, int32, int64:
			return fmt.Sprintf("₹ %d", val)
		case []uint8:
			return "₹ " + string(val)
		case pgtype.Numeric:
			if !val.Valid {
				return "₹ 0"
			}
			f, _ := val.Float64Value()
			return fmt.Sprintf("₹ %.0f", f.Float64)
		default:
			return fmt.Sprintf("₹ %v", val)
		}
	}

	return PaymentStats{
		TotalPayments:  row.TotalPayments,
		TotalReceived:  formatAmt(row.TotalReceived),
		PendingAmount:  formatAmt(row.PendingAmount),
		PendingCount:   row.PendingCount,
		FailedAmount:   formatAmt(row.FailedAmount),
		FailedCount:    row.FailedCount,
		RefundedAmount: formatAmt(row.RefundedAmount),
		RefundedCount:  row.RefundedCount,
	}, nil
}

type CreatePaymentRequest struct {
	PaymentNumber    string    `json:"payment_number"`
	TenantID         string    `json:"tenant_id"`
	SubscriptionID   string    `json:"tenant_subscription_id"`
	StatusID         string    `json:"payment_status_id"`
	ModeID           string    `json:"payment_mode_id"`
	PaymentDate      time.Time `json:"payment_date"`
	Amount           float64   `json:"amount"`
	Remarks          string    `json:"remarks"`
	CreatedBy        string    `json:"created_by"`
}

type UpdatePaymentRequest struct {
	StatusID         string    `json:"payment_status_id"`
	ModeID           string    `json:"payment_mode_id"`
	PaymentDate      *time.Time `json:"payment_date"`
	Amount           *float64   `json:"amount"`
	Remarks          string    `json:"remarks"`
	UpdatedBy        string    `json:"updated_by"`
}

func (s *QueryService) Create(ctx context.Context, req CreatePaymentRequest) (db.TenantSubscriptionPayment, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return db.TenantSubscriptionPayment{}, err
	}

	var pgID pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	var pgTenantID, pgSubID, pgStatusID, pgModeID, pgCreatedBy pgtype.UUID
	if req.TenantID != "" { tId, _ := uuid.Parse(req.TenantID); pgTenantID.Bytes = tId; pgTenantID.Valid = true }
	if req.SubscriptionID != "" { sId, _ := uuid.Parse(req.SubscriptionID); pgSubID.Bytes = sId; pgSubID.Valid = true }
	if req.StatusID != "" { stId, _ := uuid.Parse(req.StatusID); pgStatusID.Bytes = stId; pgStatusID.Valid = true }
	if req.ModeID != "" { mId, _ := uuid.Parse(req.ModeID); pgModeID.Bytes = mId; pgModeID.Valid = true }
	if req.CreatedBy != "" { cId, _ := uuid.Parse(req.CreatedBy); pgCreatedBy.Bytes = cId; pgCreatedBy.Valid = true }

	var pgAmount pgtype.Numeric
	pgAmount.Scan(fmt.Sprintf("%f", req.Amount))

	var pgDate pgtype.Date
	pgDate.Time = req.PaymentDate
	if req.PaymentDate.IsZero() { pgDate.Time = time.Now() }
	pgDate.Valid = true

	var pgRemarks pgtype.Text
	if req.Remarks != "" {
		pgRemarks.String = req.Remarks
		pgRemarks.Valid = true
	}

	var pgNow pgtype.Timestamptz
	pgNow.Time = time.Now()
	pgNow.Valid = true

	return s.repo.Create(ctx, db.CreatePaymentParams{
		TenantSubscriptionPaymentUuid: pgID,
		PaymentNumber:                req.PaymentNumber,
		TenantUuid:                   pgTenantID,
		TenantSubscriptionUuid:       pgSubID,
		PaymentStatusUniversalUuid:   pgStatusID,
		PaymentModeUniversalUuid:     pgModeID,
		PaymentDate:                  pgDate,
		Amount:                       pgAmount,
		Remarks:                      pgRemarks,
		CreatedAt:                    pgNow,
		CreatedByUserUuid:            pgCreatedBy,
	})
}

func (s *QueryService) Get(ctx context.Context, id uuid.UUID) (db.TenantSubscriptionPayment, error) {
	return s.repo.Get(ctx, id)
}

func (s *QueryService) Update(ctx context.Context, id uuid.UUID, req UpdatePaymentRequest) (db.TenantSubscriptionPayment, error) {
	var pgID pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	var pgStatusID, pgModeID, pgUpdatedBy pgtype.UUID
	if req.StatusID != "" { stId, _ := uuid.Parse(req.StatusID); pgStatusID.Bytes = stId; pgStatusID.Valid = true }
	if req.ModeID != "" { mId, _ := uuid.Parse(req.ModeID); pgModeID.Bytes = mId; pgModeID.Valid = true }
	if req.UpdatedBy != "" { uId, _ := uuid.Parse(req.UpdatedBy); pgUpdatedBy.Bytes = uId; pgUpdatedBy.Valid = true }

	var pgAmount pgtype.Numeric
	if req.Amount != nil {
		pgAmount.Scan(fmt.Sprintf("%f", *req.Amount))
	}

	var pgDate pgtype.Date
	if req.PaymentDate != nil {
		pgDate.Time = *req.PaymentDate
		pgDate.Valid = true
	}

	var pgRemarks pgtype.Text
	if req.Remarks != "" {
		pgRemarks.String = req.Remarks
		pgRemarks.Valid = true
	}

	var pgNow pgtype.Timestamptz
	pgNow.Time = time.Now()
	pgNow.Valid = true

	return s.repo.Update(ctx, db.UpdatePaymentParams{
		TenantSubscriptionPaymentUuid: pgID,
		PaymentStatusUniversalUuid:    pgStatusID,
		PaymentModeUniversalUuid:      pgModeID,
		PaymentDate:                   pgDate,
		Amount:                        pgAmount,
		Remarks:                       pgRemarks,
		UpdatedAt:                     pgNow,
		UpdatedByUserUuid:             pgUpdatedBy,
	})
}

func (s *QueryService) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	var pgID, pgDeletedBy pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	if deletedBy != "" { dId, _ := uuid.Parse(deletedBy); pgDeletedBy.Bytes = dId; pgDeletedBy.Valid = true }

	var pgNow pgtype.Timestamptz
	pgNow.Time = time.Now()
	pgNow.Valid = true

	return s.repo.Delete(ctx, db.DeletePaymentParams{
		TenantSubscriptionPaymentUuid: pgID,
		DeletedAt:                     pgNow,
		DeletedByUserUuid:             pgDeletedBy,
	})
}
