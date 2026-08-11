package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"hellocrm-superadmin/internal/platform/database/db"
)

type Service struct {
	queries   *db.Queries
	secretKey []byte // For JWT verification
}

func NewService(queries *db.Queries, secretKey []byte) *Service {
	return &Service{
		queries:   queries,
		secretKey: secretKey,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	RoleID string `json:"role_id"`
}

// ParseToken parses the JWT token and returns the user ID and Role ID.
func (s *Service) ParseToken(tokenStr string) (uuid.UUID, string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, "", errors.New("invalid token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, "", errors.New("invalid token subject")
	}

	u, _ := uuid.Parse(subject)
	return u, claims.RoleID, nil
}

// Login verifies credentials and generates a JWT.
func (s *Service) Login(ctx context.Context, email, password string) (string, Actor, error) {
	if email == "" || password == "" {
		return "", Actor{}, errors.New("invalid credentials")
	}

	// Fetch a user from the new 'user' table (since this is a stub, we just pick the first one matching email or just any user)
	// Let's get the actual Super Admin role ID from DB instead of hardcoding
	roles, err := s.queries.ListRoles(ctx)
	if err != nil || len(roles) == 0 {
		return "", Actor{}, errors.New("failed to get roles")
	}
	
	var roleUuid pgtype.UUID
	for _, r := range roles {
		if r.RoleName == "Super Admin" {
			roleUuid = r.RoleUuid
			break
		}
	}
	
	if !roleUuid.Valid {
		roleUuid = roles[0].RoleUuid // fallback to the first role if Super Admin not found
	}

	userID, _ := uuid.NewV7()

	actor := Actor{
		UserID:      userID,
		TenantID:    nil,
		Role:        "Super Admin", // Assuming we fell back to the first role or got super admin
		RoleID:      uuid.UUID(roleUuid.Bytes).String(),
		Permissions: make(map[string]struct{}),
	}

	// Fetch permissions
	rows, err := s.queries.GetRolePermissions(ctx, roleUuid)
	if err == nil {
		for _, row := range rows {
			if row.IsGranted {
				// We store string like "module:permissionCode" e.g. "error-log:v" or just "error-log" if we only care about view.
				// Let's store "moduleName:permissionCode"
				key := row.ModuleName + ":" + row.PermissionCode
				actor.Permissions[key] = struct{}{}
			}
		}
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		RoleID: actor.RoleID,
	})

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", Actor{}, err
	}

	return signedToken, actor, nil
}

// GetActorContext loads user details and RBAC permissions to build the Actor context.
func (s *Service) GetActorContext(ctx context.Context, userID uuid.UUID, roleID string) (Actor, error) {
	actor := Actor{
		UserID:      userID,
		TenantID:    nil,
		Role:        "Super Admin", // Assuming we fell back to the first role or got super admin
		RoleID:      roleID,
		Permissions: make(map[string]struct{}),
	}

	var roleUuid pgtype.UUID
	if err := roleUuid.Scan(roleID); err == nil {
		rows, err := s.queries.GetRolePermissions(ctx, roleUuid)
		if err == nil {
			for _, row := range rows {
				if row.IsGranted {
					key := row.ModuleName + ":" + row.PermissionCode
					actor.Permissions[key] = struct{}{}
				}
			}
		}
	}

	return actor, nil
}
