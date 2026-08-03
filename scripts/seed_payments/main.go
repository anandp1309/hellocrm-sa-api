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
)

func main() {
	_ = godotenv.Overload("../../.env")
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	fmt.Println("Seeding Addon Type & Addon Category...")

	// 1. Insert Master Types (Addon Type, Addon Category)
	addonTypeID, _ := uuid.NewV7()
	_, _ = conn.Exec(ctx, "INSERT INTO mst_type (type_uuid, type_name, display_order, is_system, is_deleted) VALUES ($1, 'Addon Type', 1, true, false) ON CONFLICT (type_name) DO NOTHING", addonTypeID)

	addonCatID, _ := uuid.NewV7()
	_, _ = conn.Exec(ctx, "INSERT INTO mst_type (type_uuid, type_name, display_order, is_system, is_deleted) VALUES ($1, 'Addon Category', 2, true, false) ON CONFLICT (type_name) DO NOTHING", addonCatID)

	fmt.Println("Master Types seeded!")

	// 2. Fetch some tenants to associate payments
	rows, err := conn.Query(ctx, "SELECT tenant_uuid FROM tenant LIMIT 10")
	var tenantIDs []uuid.UUID
	if err == nil {
		for rows.Next() {
			var id uuid.UUID
			_ = rows.Scan(&id)
			tenantIDs = append(tenantIDs, id)
		}
		rows.Close()
	}

	if len(tenantIDs) > 0 {
		fmt.Println("Seeding Mock Customer Payments...")
		
		// Ensure a Universal Master exists for 'Paid' payment status
		statusID, _ := uuid.NewV7()
		var dbStatusID uuid.UUID
		err = conn.QueryRow(ctx, "SELECT universal_uuid FROM mst_universal WHERE value_name = 'Paid' LIMIT 1").Scan(&dbStatusID)
		if err != nil {
			// Create it if it doesn't exist
			conn.Exec(ctx, "INSERT INTO mst_universal (universal_uuid, value_name, is_deleted) VALUES ($1, 'Paid', false)", statusID)
			dbStatusID = statusID
		}

		modeID, _ := uuid.NewV7()
		var dbModeID uuid.UUID
		err = conn.QueryRow(ctx, "SELECT universal_uuid FROM mst_universal WHERE value_name = 'VISA' LIMIT 1").Scan(&dbModeID)
		if err != nil {
			conn.Exec(ctx, "INSERT INTO mst_universal (universal_uuid, value_name, is_deleted) VALUES ($1, 'VISA', false)", modeID)
			dbModeID = modeID
		}

		// Seed Payments
		for i := 0; i < 20; i++ {
			payID, _ := uuid.NewV7()
			payNumber := fmt.Sprintf("PAY-2026-%04d", 100+i)
			randTenant := tenantIDs[rand.Intn(len(tenantIDs))]
			amt := rand.Float64()*10000 + 1000

			_, err = conn.Exec(ctx, `
				INSERT INTO tenant_subscription_payment 
				(tenant_subscription_payment_uuid, payment_number, tenant_uuid, payment_status_universal_uuid, payment_mode_universal_uuid, payment_date, amount)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT DO NOTHING
			`, payID, payNumber, randTenant, dbStatusID, dbModeID, time.Now().Add(-time.Hour * time.Duration(rand.Intn(200))), amt)
			
			if err != nil {
				log.Printf("Error inserting payment: %v", err)
			}
		}
		fmt.Println("Payments seeded successfully!")
	} else {
		fmt.Println("No tenants found! Please make sure tenants exist before seeding payments.")
	}
}
