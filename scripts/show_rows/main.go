package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Overload("C:/hellocrm-superadmin/EZ_Engineering_OS_v1.0_Final/.env")
	
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// Show Tenants
	fmt.Println("=== RECENTLY INSERTED TENANTS ===")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED_AT")
	
	rows, err := conn.Query(ctx, "SELECT id, name, created_at FROM tenants ORDER BY created_at DESC LIMIT 5")
	if err == nil {
		for rows.Next() {
			var id, name, createdAt string
			rows.Scan(&id, &name, &createdAt)
			fmt.Fprintf(w, "%s\t%s\t%s\n", id, name, createdAt)
		}
		rows.Close()
	}
	w.Flush()

	fmt.Println("\n=== RECENTLY INSERTED USERS ===")
	w2 := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w2, "ID\tEMAIL\tTENANT_ID\tCREATED_AT")
	
	rows2, err := conn.Query(ctx, "SELECT id, email, COALESCE(tenant_id::text, 'NULL'), created_at FROM users ORDER BY created_at DESC LIMIT 10")
	if err == nil {
		for rows2.Next() {
			var id, email, tenantId, createdAt string
			rows2.Scan(&id, &email, &tenantId, &createdAt)
			fmt.Fprintf(w2, "%s\t%s\t%s\t%s\n", id, email, tenantId, createdAt)
		}
		rows2.Close()
	}
	w2.Flush()
}
