package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env (checking current dir and up one level)
	// Force load the exact absolute path to the API project's .env file
	err := godotenv.Overload("C:/hellocrm-superadmin/EZ_Engineering_OS_v1.0_Final/.env")
	if err != nil {
		log.Println("Could not load absolute path .env file, continuing anyway...")
	}
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in environment or .env file")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	fmt.Println("Starting massive database seeding (50+ tenants)...")
	
	// Create common password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// 1. Ensure Roles exist
	roleNames := []string{"Super Admin", "Tenant Admin", "User"}
	roleIDs := make(map[string]uuid.UUID)

	for _, name := range roleNames {
		id, _ := uuid.NewV7()
		_, err = conn.Exec(ctx, "INSERT INTO roles (id, name, description) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING", id, name, name+" role")
		if err != nil {
			log.Fatalf("Failed to insert role %s: %v", name, err)
		}
		
		var dbID uuid.UUID
		err = conn.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", name).Scan(&dbID)
		if err != nil {
			log.Fatalf("Failed to retrieve role ID for %s: %v", name, err)
		}
		roleIDs[name] = dbID
	}

	// 2. Ensure Plans exist
	planData := []struct{
		Name string
		Price float64
	}{
		{"Basic Plan", 29.99},
		{"Pro Plan", 99.99},
		{"Enterprise Plan", 299.99},
	}
	planIDs := make([]uuid.UUID, 0)
	
	for _, p := range planData {
		id, _ := uuid.NewV7()
		_, err = conn.Exec(ctx, "INSERT INTO plans (id, name, description, price, interval) VALUES ($1, $2, $3, $4, 'monthly') ON CONFLICT DO NOTHING", id, p.Name, p.Name+" description", p.Price)
		if err != nil {
			log.Fatalf("Failed to insert plan %s: %v", p.Name, err)
		}
		
		var dbID uuid.UUID
		err = conn.QueryRow(ctx, "SELECT id FROM plans WHERE name = $1", p.Name).Scan(&dbID)
		if err != nil {
			log.Fatalf("Failed to retrieve plan ID for %s: %v", p.Name, err)
		}
		planIDs = append(planIDs, dbID)
	}

	// 3. Generate 50 Tenants, Subscriptions, and Users
	rand.Seed(time.Now().UnixNano())
	
	for i := 400; i < 450; i++ {
		tenantID, _ := uuid.NewV7()
		tenantName := fmt.Sprintf("Company %d LLC", i)
		
		// Insert Tenant
		_, err = conn.Exec(ctx, "INSERT INTO tenants (id, name) VALUES ($1, $2)", tenantID, tenantName)
		if err != nil {
			log.Printf("Warning: Failed to insert tenant %s: %v", tenantName, err)
			continue
		}

		// Insert Subscription (Random plan)
		subID, _ := uuid.NewV7()
		randomPlanID := planIDs[rand.Intn(len(planIDs))]
		_, err = conn.Exec(ctx, "INSERT INTO subscriptions (id, tenant_id, plan_id, status) VALUES ($1, $2, $3, 'active')", subID, tenantID, randomPlanID)
		if err != nil {
			log.Printf("Warning: Failed to insert subscription for %s: %v", tenantName, err)
		}

		// Insert 1 Tenant Admin
		adminID, _ := uuid.NewV7()
		adminEmail := fmt.Sprintf("admin@company%d.com", i)
		_, err = conn.Exec(ctx, "INSERT INTO users (id, email, password_hash, tenant_id, role_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
			adminID, adminEmail, string(passwordHash), tenantID, roleIDs["Tenant Admin"])
		
		// Insert 10 to 30 Users
		numUsers := rand.Intn(21) + 10
		for j := 1; j <= numUsers; j++ {
			userID, _ := uuid.NewV7()
			userEmail := fmt.Sprintf("user%d@company%d.com", j, i)
			_, err = conn.Exec(ctx, "INSERT INTO users (id, email, password_hash, tenant_id, role_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
				userID, userEmail, string(passwordHash), tenantID, roleIDs["User"])
		}
	}
	
	// Create global super admin if not exists
	superAdminID, _ := uuid.NewV7()
	_, err = conn.Exec(ctx, "INSERT INTO users (id, email, password_hash, tenant_id, role_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		superAdminID, "superadmin@hellocrm.com", string(passwordHash), nil, roleIDs["Super Admin"])

	// Lastly, just for safety, update ALL users in the DB to "password123"
	_, err = conn.Exec(ctx, "UPDATE users SET password_hash = $1", string(passwordHash))
	if err != nil {
		log.Fatalf("Failed to bulk update passwords: %v", err)
	}

	fmt.Println("Successfully seeded 50+ tenants and their users!")
	fmt.Println("All users in the database now have the password: password123")
	fmt.Println("Use superadmin@hellocrm.com for global access.")
}
