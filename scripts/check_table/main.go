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
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	tables := []string{"mst_currency", "mst_status"}
	for _, t := range tables {
		rows, err := conn.Query(ctx, "SELECT column_name FROM information_schema.columns WHERE table_name = $1", t)
		if err != nil {
			continue
		}
		
		fmt.Printf("\n--- Table: %s ---\n", t)
		for rows.Next() {
			var col string
			rows.Scan(&col)
			fmt.Println(col)
		}
		rows.Close()
	}
}
