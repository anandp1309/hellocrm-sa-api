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

// ParseToken parses the JWT token and returns the user ID.
func (s *Service) ParseToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, errors.New("invalid token subject")
	}

	u, _ := uuid.Parse(subject)
	return u, nil
}

// Login verifies credentials and generates a JWT.
func (s *Service) Login(ctx context.Context, email, password string) (string, Actor, error) {
	// TODO: Fetch user by email from the database
	// user, err := s.queries.GetUserByEmail(ctx, email)
	// if err != nil { return "", Actor{}, err }

	// TODO: Verify password using bcrypt
	// if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
	//     return "", Actor{}, errors.New("invalid credentials")
	// }

	// STUB: Accepting any non-empty credentials for now
	if email == "" || password == "" {
		return "", Actor{}, errors.New("invalid credentials")
	}

	// STUB: Generate a dummy user ID for the token
	userID, _ := uuid.NewV7()

	// Fetch Role and Permissions using the same logic we use for middleware
	actor, err := s.GetActorContext(ctx, userID)
	if err != nil {
		return "", Actor{}, err
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", Actor{}, err
	}

	return signedToken, actor, nil
}

// GetActorContext loads user details and RBAC permissions to build the Actor context.
func (s *Service) GetActorContext(ctx context.Context, userID uuid.UUID) (Actor, error) {
	// 1. In a complete flow, we might fetch the user record first to get the TenantID.
	// user, err := s.queries.GetUserByID(ctx, userID)

	// 2. Load Permissions mapping for this user's role.
	perms, err := s.queries.GetUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return Actor{}, errors.New("failed to fetch user permissions")
	}

	permMap := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		permMap[p] = struct{}{}
	}

	return Actor{
		UserID:      userID,
		TenantID:    nil, // Would be user.TenantID
		Role:        "fetched_role_name",
		RoleID:      "super_admin",
		Permissions: permMap,
	}, nil
}
