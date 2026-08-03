package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Overload("C:/hellocrm-superadmin/EZ_Engineering_OS_v1.0_Final/.env")
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	fmt.Printf("=== DB CONNECTION INFO ===\n")
	fmt.Printf("URL: %s\n", dsn)
	fmt.Println("==========================")

	var tenantCount int
	err = conn.QueryRow(ctx, "SELECT count(*) FROM tenants").Scan(&tenantCount)
	var userCount int
	err = conn.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&userCount)
	
	fmt.Printf("Database actually contains: %d tenants and %d users.\n", tenantCount, userCount)

	fmt.Println("\n--- SAMPLE USERS (First 5) ---")
	rows, err := conn.Query(ctx, "SELECT email FROM users LIMIT 5")
	if err == nil {
		for rows.Next() {
			var email string
			rows.Scan(&email)
			fmt.Println(email)
		}
		rows.Close()
	}

	fmt.Println("\n--- SAMPLE TENANTS (First 5) ---")
	rows, err = conn.Query(ctx, "SELECT name FROM tenants LIMIT 5")
	if err == nil {
		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Println(name)
		}
		rows.Close()
	}
}
