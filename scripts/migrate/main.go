package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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

	migrationsDir := filepath.Join("db", "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			log.Fatalf("Failed to read migration %s: %v", file, err)
		}
		
		fmt.Printf("Applying migration: %s\n", file)
		_, err = conn.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}

	fmt.Println("Migrations applied successfully!")
}
