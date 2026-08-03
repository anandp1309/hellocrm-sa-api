package role

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type QueryService interface {
	ListRoles(ctx context.Context) ([]RoleView, error)
	GetRolePermissions(ctx context.Context, id string) ([]PermissionView, error)
}

type queryService struct {
	repo Repository
}

func NewQueryService(repo Repository) QueryService {
	return &queryService{repo: repo}
}

type RoleView struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Count      int64  `json:"count"`
	Active     bool   `json:"active"`
	ColorClass string `json:"colorClass"`
	Icon       string `json:"icon"`
}

type PermissionView struct {
	Module    string `json:"module"`
	IconColor string `json:"iconColor"`
	V         string `json:"v"`
	C         string `json:"c"`
	E         string `json:"e"`
	D         string `json:"d"`
	X         string `json:"x"`
	I         string `json:"i"`
	M         string `json:"m"`
}

func (s *queryService) ListRoles(ctx context.Context) ([]RoleView, error) {
	rows, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	var data []RoleView
	for i, row := range rows {
		colorClass := "bg-gray-100 text-gray-600 border-transparent"
		icon := "M15 12a3 3 0 11-6 0 3 3 0 016 0zM2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
		
		nameLower := strings.ToLower(row.RoleName)
		if strings.Contains(nameLower, "super admin") {
			colorClass = "bg-indigo-50 text-indigo-600 border-indigo-100"
			icon = "M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
		} else if strings.Contains(nameLower, "admin") {
			colorClass = "bg-blue-50 text-blue-600 border-transparent"
			icon = "M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
		} else if strings.Contains(nameLower, "manager") {
			colorClass = "bg-green-50 text-green-600 border-transparent"
			icon = "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
		} else if strings.Contains(nameLower, "support") {
			colorClass = "bg-orange-50 text-orange-600 border-transparent"
			icon = "M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a4.978 4.978 0 01-1.414-2.83m-1.414 5.658a9 9 0 01-2.167-9.238m7.824 2.168a2 2 0 11-2.829-2.83 2 2 0 012.829 2.83z"
		}

		uuidStr := fmt.Sprintf("%x-%x-%x-%x-%x", row.RoleUuid.Bytes[0:4], row.RoleUuid.Bytes[4:6], row.RoleUuid.Bytes[6:8], row.RoleUuid.Bytes[8:10], row.RoleUuid.Bytes[10:16])
		
		data = append(data, RoleView{
			ID:         uuidStr,
			Name:       row.RoleName,
			Desc:       row.Remarks.String,
			Count:      row.UserCount,
			Active:     i == 0,
			ColorClass: colorClass,
			Icon:       icon,
		})
	}
	return data, nil
}

func (s *queryService) GetRolePermissions(ctx context.Context, id string) ([]PermissionView, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.GetRolePermissions(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Map of module -> permission map
	moduleMap := make(map[string]map[string]string)
	
	for _, row := range rows {
		code := row.PermissionCode
		parts := strings.Split(code, "_")
		if len(parts) >= 2 {
			module := parts[0]
			action := parts[1]
			
			if _, ok := moduleMap[module]; !ok {
				moduleMap[module] = map[string]string{
					"v": "n", "c": "n", "e": "n", "d": "n", "x": "n", "i": "n", "m": "n",
				}
			}
			
			actionLower := strings.ToLower(action)
			if actionLower == "view" { moduleMap[module]["v"] = "y" }
			if actionLower == "create" { moduleMap[module]["c"] = "y" }
			if actionLower == "edit" { moduleMap[module]["e"] = "y" }
			if actionLower == "delete" { moduleMap[module]["d"] = "y" }
			if actionLower == "export" { moduleMap[module]["x"] = "y" }
			if actionLower == "import" { moduleMap[module]["i"] = "y" }
			if actionLower == "manage" { moduleMap[module]["m"] = "y" }
		}
	}

	// Default modules if none exist (for UI testing)
	if len(moduleMap) == 0 {
		return []PermissionView{
			{Module: "Dashboard", IconColor: "text-indigo-600", V: "y", C: "y", E: "y", D: "y", X: "y", I: "y", M: "y"},
			{Module: "Customers", IconColor: "text-blue-500", V: "y", C: "y", E: "y", D: "y", X: "y", I: "y", M: "y"},
			{Module: "Subscriptions", IconColor: "text-green-500", V: "y", C: "y", E: "y", D: "y", X: "y", I: "y", M: "y"},
			{Module: "Roles & Permissions", IconColor: "text-green-600", V: "y", C: "y", E: "y", D: "y", X: "y", I: "y", M: "y"},
		}, nil
	}

	var perms []PermissionView
	for module, actions := range moduleMap {
		perms = append(perms, PermissionView{
			Module:    module,
			IconColor: "text-indigo-600",
			V:         actions["v"],
			C:         actions["c"],
			E:         actions["e"],
			D:         actions["d"],
			X:         actions["x"],
			I:         actions["i"],
			M:         actions["m"],
		})
	}

	return perms, nil
}
