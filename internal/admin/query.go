package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryService struct {
	repo *Repository
}

func NewQueryService(repo *Repository) *QueryService {
	return &QueryService{repo: repo}
}

type AdminView struct {
	ID        string `json:"id"`
	Initials  string `json:"initials"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Mfa       string `json:"mfa"`
	LastLogin string `json:"lastLogin"`
	Color     string `json:"color"`
}

type PaginatedAdmins struct {
	Data       []AdminView `json:"data"`
	Total      int64       `json:"total"`
	TotalPages int32       `json:"totalPages"`
}

func getInitials(firstName, lastName string) string {
	initials := ""
	if len(firstName) > 0 {
		initials += string(firstName[0])
	}
	if len(lastName) > 0 {
		initials += string(lastName[0])
	}
	return strings.ToUpper(initials)
}

func (s *QueryService) ListPaginated(ctx context.Context, search, role, status, mfa string, page, limit int32) (PaginatedAdmins, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := s.repo.ListPaginated(ctx, search, role, status, mfa, limit, offset)
	if err != nil {
		return PaginatedAdmins{}, err
	}

	var data []AdminView
	var total int64
	for i, row := range rows {
		total = row.TotalRecords
		
		lastLogin := "—"
		if row.LastLoginAt.Valid {
			lastLogin = row.LastLoginAt.Time.Format("02 Jan 2006, 03:04 PM")
		}

		color := "bg-indigo-700 text-white"
		if i%3 == 1 {
			color = "bg-blue-600 text-white"
		} else if i%3 == 2 {
			color = "bg-teal-600 text-white"
		}

		name := row.FirstName
		if row.LastName.Valid {
			name += " " + row.LastName.String
		}

		data = append(data, AdminView{
			ID:        fmt.Sprintf("%x-%x-%x-%x-%x", row.UserUuid.Bytes[0:4], row.UserUuid.Bytes[4:6], row.UserUuid.Bytes[6:8], row.UserUuid.Bytes[8:10], row.UserUuid.Bytes[10:16]),
			Initials:  getInitials(row.FirstName, row.LastName.String),
			Name:      name,
			Username:  row.Username,
			Email:     row.EmailAddress.String,
			Role:      row.Role.String,
			Status:    row.Status,
			Mfa:       row.Mfa,
			LastLogin: lastLogin,
			Color:     color,
		})
	}
	
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return PaginatedAdmins{
		Data:       data,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

type AdminStats struct {
	TotalUsers    int64 `json:"totalUsers"`
	ActiveUsers   int64 `json:"activeUsers"`
	InactiveUsers int64 `json:"inactiveUsers"`
	TotalRoles    int64 `json:"totalRoles"`
}

func (s *QueryService) GetStats(ctx context.Context) (AdminStats, error) {
	row, err := s.repo.GetStats(ctx)
	if err != nil {
		return AdminStats{}, err
	}
	return AdminStats{
		TotalUsers:    row.TotalUsers,
		ActiveUsers:   row.ActiveUsers,
		InactiveUsers: row.InactiveUsers,
		TotalRoles:    row.TotalRoles,
	}, nil
}

func (s *QueryService) GetAdmin(ctx context.Context, id uuid.UUID) (interface{}, error) {
	var pgUserId pgtype.UUID
	pgUserId.Bytes = id
	pgUserId.Valid = true
	return s.repo.q.GetAdmin(ctx, pgUserId)
}
