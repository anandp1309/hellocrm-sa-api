package auth

import (
	"context"
	"fmt"
	"log"
	"strings"

	"hellocrm-superadmin/internal/platform/database/db"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

type CasbinEnforcer struct {
	Enforcer *casbin.Enforcer
	Queries  *db.Queries
}

func NewCasbinEnforcer(queries *db.Queries) (*CasbinEnforcer, error) {
	// Initialize Casbin model from string
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "super_admin"
`)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	ce := &CasbinEnforcer{
		Enforcer: e,
		Queries:  queries,
	}

	if err := ce.LoadPolicies(context.Background()); err != nil {
		log.Printf("Failed to load Casbin policies: %v", err)
	}

	return ce, nil
}

func (ce *CasbinEnforcer) LoadPolicies(ctx context.Context) error {
	perms, err := ce.Queries.GetAllRolePermissions(ctx)
	if err != nil {
		return err
	}

	ce.Enforcer.ClearPolicy()

	for _, p := range perms {
		// p.RoleUuid is pgtype.UUID
		roleID := fmt.Sprintf("%x-%x-%x-%x-%x", p.RoleUuid.Bytes[0:4], p.RoleUuid.Bytes[4:6], p.RoleUuid.Bytes[6:8], p.RoleUuid.Bytes[8:10], p.RoleUuid.Bytes[10:16])
		
		code := p.PermissionCode // e.g. DASHBOARD_VIEW
		parts := strings.Split(code, "_")
		if len(parts) >= 2 {
			obj := strings.ToLower(parts[0]) // dashboard
			act := strings.ToLower(parts[1]) // view
			
			_, err := ce.Enforcer.AddPolicy(roleID, obj, act)
			if err != nil {
				log.Printf("Error adding policy: %v", err)
			}
		}
	}

	return nil
}

func (ce *CasbinEnforcer) Enforce(roleID, obj, act string) (bool, error) {
	return ce.Enforcer.Enforce(roleID, obj, act)
}
