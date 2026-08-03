package universal

import (
	"context"
	"sort"
	"strings"
	"time"

	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateUniversalRequest struct {
	MasterTypeID string `json:"master_type_id"`
	ValueName    string `json:"value_name"`
	DisplayOrder int32  `json:"display_order"`
	IsSystem     bool   `json:"is_system"`
	Remarks      string `json:"remarks"`
	CreatedBy    string `json:"created_by"`
}

func (s *Service) Create(ctx context.Context, req CreateUniversalRequest) (db.MstUniversal, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return db.MstUniversal{}, err
	}

	masterTypeUUID, err := uuid.Parse(req.MasterTypeID)
	if err != nil {
		return db.MstUniversal{}, err
	}

	createdByUser, _ := uuid.Parse(req.CreatedBy)

	var pgID, pgMasterTypeID, pgCreatedBy pgtype.UUID
	pgID.Bytes = id
	pgID.Valid = true

	pgMasterTypeID.Bytes = masterTypeUUID
	pgMasterTypeID.Valid = true

	if createdByUser != uuid.Nil {
		pgCreatedBy.Bytes = createdByUser
		pgCreatedBy.Valid = true
	}

	var pgDisplayOrder pgtype.Int4
	pgDisplayOrder.Int32 = req.DisplayOrder
	pgDisplayOrder.Valid = true

	var pgIsSystem pgtype.Bool
	pgIsSystem.Bool = req.IsSystem
	pgIsSystem.Valid = true

	var pgRemarks pgtype.Text
	if req.Remarks != "" {
		pgRemarks.String = req.Remarks
		pgRemarks.Valid = true
	}

	var pgNow pgtype.Timestamptz
	pgNow.Time = time.Now()
	pgNow.Valid = true

	return s.repo.Create(ctx, db.CreateUniversalParams{
		UniversalUuid:      pgID,
		MasterTypeUuid:     pgMasterTypeID,
		ValueName:          req.ValueName,
		DisplayOrder:       req.DisplayOrder,
		IsSystem:           req.IsSystem,
		Remarks:            pgRemarks,
		CreatedAt:          pgNow,
		CreatedByUserUuid:  pgCreatedBy,
	})
}

func (s *Service) Update(ctx context.Context, id string, req CreateUniversalRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	
	pgRemarks := pgtype.Text{}
	if req.Remarks != "" {
		pgRemarks.String = req.Remarks
		pgRemarks.Valid = true
	}
	
	return s.repo.Update(ctx, uid, req.ValueName, pgRemarks)
}

func (s *Service) Deactivate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Deactivate(ctx, uid)
}

type PaginatedUniversal struct {
	Data           []db.MstUniversal `json:"data"`
	TotalRecords   int               `json:"total_records"`
	PresentRecords int               `json:"present_records"`
	CurrentPage    int               `json:"current_page"`
	TotalPages     int               `json:"total_pages"`
	Limit          int               `json:"limit"`
}

type UniversalFilter struct {
	Search    string
	Scope     string
	Tenant    string
	Status    string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

func (s *Service) GetByMasterType(ctx context.Context, masterTypeID string, filter UniversalFilter) (PaginatedUniversal, error) {
	id, err := uuid.Parse(masterTypeID)
	if err != nil {
		return PaginatedUniversal{}, err
	}
	
	dbVals, err := s.repo.GetByMasterType(ctx, id)
	if err != nil {
		return PaginatedUniversal{}, err
	}

	var allVals []db.MstUniversal
	for _, v := range dbVals {
		// Filter by search
		if filter.Search != "" && !strings.Contains(strings.ToLower(v.ValueName), strings.ToLower(filter.Search)) {
			continue
		}
		
		// Filter by status (is_deleted)
		if filter.Status != "" && filter.Status != "All Statuses" {
			isActive := !v.IsDeleted
			if filter.Status == "Active" && !isActive {
				continue
			}
			if filter.Status == "Inactive" && isActive {
				continue
			}
		}
		
		// Note: Scope and Tenant are currently hardcoded as Global and — in UI
		// Filtering for them would be mocked here
		if filter.Scope != "" && filter.Scope != "All" && filter.Scope != "Global" {
			continue
		}
		
		allVals = append(allVals, v)
	}

	totalRecords := len(allVals)

	// Sorting
	if filter.SortBy != "" {
		sort.Slice(allVals, func(i, j int) bool {
			var isLess bool
			switch filter.SortBy {
			case "value":
				isLess = strings.ToLower(allVals[i].ValueName) < strings.ToLower(allVals[j].ValueName)
			case "status":
				isLess = (!allVals[i].IsDeleted) && (allVals[j].IsDeleted)
			case "createdOn":
				isLess = allVals[i].CreatedAt.Time.Before(allVals[j].CreatedAt.Time)
			default:
				isLess = strings.ToLower(allVals[i].ValueName) < strings.ToLower(allVals[j].ValueName)
			}
			if filter.SortOrder == "desc" {
				return !isLess
			}
			return isLess
		})
	}

	// Pagination
	var pagedVals []db.MstUniversal
	if filter.Limit == -1 {
		pagedVals = allVals
	} else {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		limit := filter.Limit
		if limit < 1 {
			limit = 10
		}
		start := (page - 1) * limit
		end := start + limit

		if start < totalRecords {
			if end > totalRecords {
				end = totalRecords
			}
			pagedVals = allVals[start:end]
		} else {
			pagedVals = []db.MstUniversal{}
		}
	}

	totalPages := 1
	if filter.Limit > 0 {
		totalPages = (totalRecords + filter.Limit - 1) / filter.Limit
	}
	if totalPages == 0 {
		totalPages = 1
	}

	if pagedVals == nil {
		pagedVals = []db.MstUniversal{}
	}

	return PaginatedUniversal{
		Data:           pagedVals,
		TotalRecords:   totalRecords,
		PresentRecords: len(pagedVals),
		CurrentPage:    filter.Page,
		TotalPages:     totalPages,
		Limit:          filter.Limit,
	}, nil
}
