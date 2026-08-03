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
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	sqlBytes, err := os.ReadFile("db/migrations/000002_add_mst_addon.up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migration successful")
}
