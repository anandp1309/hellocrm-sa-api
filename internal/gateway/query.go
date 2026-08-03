package gateway

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

type GatewayTransactionView struct {
	ID             string `json:"id"`
	TransactionId  string `json:"transactionId"`
	Gateway        string `json:"gateway"`
	Customer       string `json:"customer"`
	Email          string `json:"email"`
	PaymentFor     string `json:"paymentFor"`
	Amount         string `json:"amount"`
	GatewayStatus  string `json:"gatewayStatus"`
	WebhookStatus  string `json:"webhookStatus"`
	PaymentStatus  string `json:"paymentStatus"`
	CreatedOn      string `json:"createdOn"`
	CreatedOnTime  string `json:"createdOnTime"`
}

type PaginatedGatewayTransactions struct {
	Data       []GatewayTransactionView `json:"data"`
	Total      int64                    `json:"total"`
	TotalPages int32                    `json:"totalPages"`
}

func (s *QueryService) ListPaginated(ctx context.Context, search, gateway, status, startDate, endDate string, page, limit int32) (PaginatedGatewayTransactions, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := s.repo.ListPaginated(ctx, search, gateway, status, startDate, endDate, limit, offset)
	if err != nil {
		return PaginatedGatewayTransactions{}, err
	}

	var data []GatewayTransactionView
	var total int64
	for _, row := range rows {
		total = row.TotalRecords
		
		f, _ := row.Amount.Float64Value()
		amt := fmt.Sprintf("₹ %.2f", f.Float64)
		
		gwStatus := "Success"
		whStatus := "Received"
		if row.PaymentStatus == "Failed" {
			gwStatus = "Failed"
			whStatus = "Error"
		} else if row.PaymentStatus == "Pending" {
			gwStatus = "Pending"
			whStatus = "Pending"
		}

		data = append(data, GatewayTransactionView{
			ID:            fmt.Sprintf("%x-%x-%x-%x-%x", row.TenantSubscriptionPaymentUuid.Bytes[0:4], row.TenantSubscriptionPaymentUuid.Bytes[4:6], row.TenantSubscriptionPaymentUuid.Bytes[6:8], row.TenantSubscriptionPaymentUuid.Bytes[8:10], row.TenantSubscriptionPaymentUuid.Bytes[10:16]),
			TransactionId: row.TransactionID,
			Gateway:       row.Gateway,
			Customer:      row.Customer,
			Email:         row.EmailAddress.String,
			PaymentFor:    row.PlanName,
			Amount:        amt,
			GatewayStatus: gwStatus,
			WebhookStatus: whStatus,
			PaymentStatus: row.PaymentStatus,
			CreatedOn:     row.CreatedAt.Time.Format("02 Jan 2006"),
			CreatedOnTime: row.CreatedAt.Time.Format("03:04 PM"),
		})
	}
	
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return PaginatedGatewayTransactions{
		Data:       data,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

type GatewayStats struct {
	TotalTransactions int64  `json:"totalTransactions"`
	TotalAmount       string `json:"totalAmount"`
	SuccessAmount     string `json:"successAmount"`
	SuccessCount      int64  `json:"successCount"`
	FailedAmount      string `json:"failedAmount"`
	FailedCount       int64  `json:"failedCount"`
	PendingAmount     string `json:"pendingAmount"`
	PendingCount      int64  `json:"pendingCount"`
}

func (s *QueryService) GetStats(ctx context.Context) (GatewayStats, error) {
	row, err := s.repo.GetStats(ctx)
	if err != nil {
		return GatewayStats{}, err
	}
	
	formatAmt := func(val pgtype.Numeric) string {
		if !val.Valid {
			return "₹ 0"
		}
		f, _ := val.Float64Value()
		return fmt.Sprintf("₹ %.0f", f.Float64)
	}

	return GatewayStats{
		TotalTransactions: row.TotalTransactions,
		TotalAmount:       formatAmt(row.TotalAmount.(pgtype.Numeric)),
		SuccessAmount:     formatAmt(row.SuccessAmount.(pgtype.Numeric)),
		SuccessCount:      row.SuccessCount,
		FailedAmount:      formatAmt(row.FailedAmount.(pgtype.Numeric)),
		FailedCount:       row.FailedCount,
		PendingAmount:     formatAmt(row.PendingAmount.(pgtype.Numeric)),
		PendingCount:      row.PendingCount,
	}, nil
}

func (s *QueryService) Get(ctx context.Context, id string) (GatewayTransactionView, error) {
	var pgId pgtype.UUID
	err := pgId.Scan(id)
	if err != nil {
		return GatewayTransactionView{}, err
	}

	row, err := s.repo.q.GetGatewayTransaction(ctx, pgId)
	if err != nil {
		return GatewayTransactionView{}, err
	}
	
	bytes, _ := row.TenantSubscriptionPaymentUuid.Value()
	var uuidStr string
	if bytes != nil {
		if b, ok := bytes.([]byte); ok && len(b) == 16 {
			uuidStr = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
		} else if s, ok := bytes.(string); ok {
			uuidStr = s
		}
	}
	
	formatAmt := func(val pgtype.Numeric) string {
		if !val.Valid {
			return "₹ 0"
		}
		f, _ := val.Float64Value()
		return fmt.Sprintf("₹ %.0f", f.Float64)
	}

	return GatewayTransactionView{
		ID:             uuidStr,
		TransactionId:  row.TransactionID,
		Gateway:        row.Gateway,
		Customer:       row.Customer,
		Email:          row.EmailAddress.String,
		PaymentFor:     "Subscription",
		Amount:         formatAmt(row.Amount),
		GatewayStatus:  "Success",
		WebhookStatus:  "Synced",
		PaymentStatus:  row.PaymentStatus,
		CreatedOn:      row.CreatedAt.Time.Format("02 Jan 2006"),
		CreatedOnTime:  row.CreatedAt.Time.Format("03:04 PM"),
	}, nil
}
