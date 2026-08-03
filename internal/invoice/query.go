package invoice

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

type InvoiceView struct {
	ID          string `json:"id"`
	Customer    string `json:"customer"`
	Email       string `json:"email"`
	Type        string `json:"type"`
	InvoiceDate string `json:"invoiceDate"`
	DueDate     string `json:"dueDate"`
	Amount      string `json:"amount"`
	Status      string `json:"status"`
	PaymentDate string `json:"paymentDate"`
}

type PaginatedInvoices struct {
	Data       []InvoiceView `json:"data"`
	Total      int64         `json:"total"`
	TotalPages int32         `json:"totalPages"`
}

func (s *QueryService) ListPaginated(ctx context.Context, search, status, startDate, endDate string, page, limit int32) (PaginatedInvoices, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := s.repo.ListPaginated(ctx, search, status, startDate, endDate, limit, offset)
	if err != nil {
		return PaginatedInvoices{}, err
	}

	var data []InvoiceView
	var total int64
	for _, row := range rows {
		total = row.TotalRecords
		
		f, _ := row.Amount.Float64Value()
		amt := fmt.Sprintf("₹ %.2f", f.Float64)
		
		payDate := "—"
		if row.PaymentDate.Valid {
			payDate = row.PaymentDate.Time.Format("02 Jan 2006")
		}

		invDate := "—"
		if row.InvoiceDate.Valid {
			invDate = row.InvoiceDate.Time.Format("02 Jan 2006")
		}

		dueDate := "—"
		if row.DueDate.Valid {
			dueDate = row.DueDate.Time.Format("02 Jan 2006")
		}

		data = append(data, InvoiceView{
			ID:          row.InvoiceNo,
			Customer:    row.Customer,
			Email:       row.EmailAddress.String,
			Type:        row.BillingType,
			InvoiceDate: invDate,
			DueDate:     dueDate,
			Amount:      amt,
			Status:      row.Status,
			PaymentDate: payDate,
		})
	}
	
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return PaginatedInvoices{
		Data:       data,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

type InvoiceStats struct {
	TotalInvoices int64  `json:"totalInvoices"`
	PaidAmount    string `json:"paidAmount"`
	PaidCount     int64  `json:"paidCount"`
	PendingAmount string `json:"pendingAmount"`
	PendingCount  int64  `json:"pendingCount"`
	OverdueAmount string `json:"overdueAmount"`
	OverdueCount  int64  `json:"overdueCount"`
}

func (s *QueryService) GetStats(ctx context.Context) (InvoiceStats, error) {
	row, err := s.repo.GetStats(ctx)
	if err != nil {
		return InvoiceStats{}, err
	}
	
	formatAmt := func(val pgtype.Numeric) string {
		if !val.Valid {
			return "₹ 0"
		}
		f, _ := val.Float64Value()
		return fmt.Sprintf("₹ %.0f", f.Float64)
	}

	return InvoiceStats{
		TotalInvoices: row.TotalInvoices,
		PaidAmount:    formatAmt(row.PaidAmount.(pgtype.Numeric)),
		PaidCount:     row.PaidCount,
		PendingAmount: formatAmt(row.PendingAmount.(pgtype.Numeric)),
		PendingCount:  row.PendingCount,
		OverdueAmount: formatAmt(row.OverdueAmount.(pgtype.Numeric)),
		OverdueCount:  row.OverdueCount,
	}, nil
}

func (s *QueryService) Get(ctx context.Context, id string) (InvoiceView, error) {
	var pgId pgtype.UUID
	err := pgId.Scan(id)
	if err != nil {
		return InvoiceView{}, err
	}

	row, err := s.repo.q.GetInvoice(ctx, pgId)
	if err != nil {
		return InvoiceView{}, err
	}
	
	bytes, _ := row.InvoiceUuid.Value()
	var uuidStr string
	if bytes != nil {
		if b, ok := bytes.([]byte); ok && len(b) == 16 {
			uuidStr = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
		} else if st, ok := bytes.(string); ok {
			uuidStr = st
		}
	}
	
	payDate := "—"
	if row.PaymentDate.Valid {
		payDate = row.PaymentDate.Time.Format("02 Jan 2006")
	}

	formatAmt := func(val pgtype.Numeric) string {
		if !val.Valid {
			return "₹ 0"
		}
		f, _ := val.Float64Value()
		return fmt.Sprintf("₹ %.0f", f.Float64)
	}
	
	return InvoiceView{
		ID:          uuidStr,
		Customer:    row.Customer,
		Email:       row.EmailAddress.String,
		Type:        row.BillingType,
		InvoiceDate: row.InvoiceDate.Time.Format("02 Jan 2006"),
		DueDate:     row.DueDate.Time.Format("02 Jan 2006"),
		Amount:      formatAmt(row.Amount),
		Status:      row.Status,
		PaymentDate: payDate,
	}, nil
}
