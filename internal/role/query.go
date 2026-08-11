package role

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"hellocrm-superadmin/internal/platform/database/db"
)

type QueryService interface {
	ListRoles(ctx context.Context) ([]RoleView, error)
	GetRolePermissions(ctx context.Context, id string) ([]PermissionView, error)
	GetPermissionsTemplate(ctx context.Context) ([]PermissionView, error)
}

type queryService struct {
	repo Repository
}

func NewQueryService(repo Repository) QueryService {
	return &queryService{repo: repo}
}

type RoleView struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Desc               string `json:"desc"`
	Count              int64  `json:"count"`
	Active             bool   `json:"active"`
	ColorClass         string `json:"colorClass"`
	Icon               string `json:"icon"`
	CreatedByUserUuid  string `json:"createdByUserUuid,omitempty"`
	CreatedByFirstName string `json:"createdByFirstName,omitempty"`
	CreatedByLastName  string `json:"createdByLastName,omitempty"`
	UpdatedByUserUuid  string `json:"updatedByUserUuid,omitempty"`
	UpdatedByFirstName string `json:"updatedByFirstName,omitempty"`
	UpdatedByLastName  string `json:"updatedByLastName,omitempty"`
}

type PermissionView struct {
	ModuleUuid string `json:"moduleUuid"`
	Module     string `json:"module"`
	IconColor  string `json:"iconColor"`
	V          bool   `json:"v"`
	V_Uuid     string `json:"v_uuid"`
	C          bool   `json:"c"`
	C_Uuid     string `json:"c_uuid"`
	E          bool   `json:"e"`
	E_Uuid     string `json:"e_uuid"`
	D          bool   `json:"d"`
	D_Uuid     string `json:"d_uuid"`
	X          bool   `json:"x"`
	X_Uuid     string `json:"x_uuid"`
	I          bool   `json:"i"`
	I_Uuid     string `json:"i_uuid"`
	M          bool   `json:"m"`
	M_Uuid     string `json:"m_uuid"`
}

type ActionInfo struct {
	Status bool
	Uuid   string
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

		uuidStr := uuid.UUID(row.RoleUuid.Bytes).String()

		var createdByUuid, updatedByUuid string
		if row.CreatedByUserUuid.Valid {
			createdByUuid = uuid.UUID(row.CreatedByUserUuid.Bytes).String()
		}
		if row.UpdatedByUserUuid.Valid {
			updatedByUuid = uuid.UUID(row.UpdatedByUserUuid.Bytes).String()
		}

		data = append(data, RoleView{
			ID:                 uuidStr,
			Name:               row.RoleName,
			Desc:               row.Remarks.String,
			Count:              row.UserCount,
			Active:             i == 0,
			ColorClass:         colorClass,
			Icon:               icon,
			CreatedByUserUuid:  createdByUuid,
			CreatedByFirstName: row.CreatedByFirstName.String,
			CreatedByLastName:  row.CreatedByLastName.String,
			UpdatedByUserUuid:  updatedByUuid,
			UpdatedByFirstName: row.UpdatedByFirstName.String,
			UpdatedByLastName:  row.UpdatedByLastName.String,
		})
	}
	return data, nil
}

type FlatPermissionView struct {
	RoleName           string `json:"roleName"`
	ModuleName         string `json:"moduleName"`
	AccessRight        string `json:"accessRight"`
	PermissionCode     string `json:"permissionCode"`
	RolePermissionUuid string `json:"rolePermissionUuid"`
	RoleUuid           string `json:"roleUuid"`
	ModuleUuid         string `json:"moduleUuid"`
	PermissionUuid     string `json:"permissionUuid"`
	TenantUuid         string `json:"tenantUuid"`
}

func (s *queryService) GetRolePermissions(ctx context.Context, id string) ([]PermissionView, error) {
	var pgRoleID pgtype.UUID
	err := pgRoleID.Scan(id)
	if err != nil {
		return nil, err
	}

	// 1. Fetch active permissions for the role
	activeRows, err := s.repo.GetRolePermissions(ctx, pgRoleID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch all permissions and modules for the template
	allPerms, err := s.repo.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	modules, err := s.repo.ListModules(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Build the template map
	templateMap := buildTemplateMap(allPerms, modules)

	// 4. Update the template map with active permissions
	for _, row := range activeRows {
		module := row.ModuleName
		key := row.PermissionCode
		if _, ok := templateMap[module]; ok {
			if _, ok2 := templateMap[module][key]; ok2 {
				// Mark as the boolean status of IsGranted
				info := templateMap[module][key]
				info.Status = row.IsGranted
				templateMap[module][key] = info
			}
		}
	}

	// 5. Convert the map to the PermissionView array
	return s.buildPermissionViews(ctx, templateMap)
}

func (s *queryService) GetPermissionsTemplate(ctx context.Context) ([]PermissionView, error) {
	allPerms, err := s.repo.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	modules, err := s.repo.ListModules(ctx)
	if err != nil {
		return nil, err
	}

	templateMap := buildTemplateMap(allPerms, modules)
	return s.buildPermissionViews(ctx, templateMap)
}

func buildTemplateMap(allPerms []db.ListAllPermissionsRow, modules []db.ListModulesRow) map[string]map[string]ActionInfo {
	templateMap := make(map[string]map[string]ActionInfo)
	for _, mRow := range modules {
		module := mRow.ModuleName
		templateMap[module] = map[string]ActionInfo{
			"v": {false, ""}, "c": {false, ""}, "e": {false, ""}, "d": {false, ""},
			"x": {false, ""}, "i": {false, ""}, "m": {false, ""},
		}
		
		for _, row := range allPerms {
			pUuidStr := uuid.UUID(row.PermissionUuid.Bytes).String()
			key := row.PermissionCode
			if _, ok := templateMap[module][key]; ok {
				templateMap[module][key] = ActionInfo{false, pUuidStr}
			}
		}
	}
	return templateMap
}

func (s *queryService) buildPermissionViews(ctx context.Context, templateMap map[string]map[string]ActionInfo) ([]PermissionView, error) {
	var perms []PermissionView
	modules, err := s.repo.ListModules(ctx)
	if err != nil {
		return nil, err
	}

	for _, mRow := range modules {
		m := mRow.ModuleName
		uuidStr := uuid.UUID(mRow.ModuleUuid.Bytes).String()
		if actions, ok := templateMap[m]; ok {
			perms = append(perms, PermissionView{
				ModuleUuid: uuidStr,
				Module:     m,
				IconColor:  "text-indigo-600",
				V:          actions["v"].Status, V_Uuid: actions["v"].Uuid,
				C:          actions["c"].Status, C_Uuid: actions["c"].Uuid,
				E:          actions["e"].Status, E_Uuid: actions["e"].Uuid,
				D:          actions["d"].Status, D_Uuid: actions["d"].Uuid,
				X:          actions["x"].Status, X_Uuid: actions["x"].Uuid,
				I:          actions["i"].Status, I_Uuid: actions["i"].Uuid,
				M:          actions["m"].Status, M_Uuid: actions["m"].Uuid,
			})
		} else {
			perms = append(perms, PermissionView{
				ModuleUuid: uuidStr,
				Module:     m,
				IconColor:  "text-indigo-600",
				V: false, C: false, E: false, D: false, X: false, I: false, M: false,
			})
		}
	}
	return perms, nil
}
