package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Addon struct {
	Name      string
	Category  string
	Limit     string
	Price     string
	Status    string
	IconSvg   string
	IconColor string
}

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

	addons := []Addon{
		{"Extra Users Pack", "Users", "5 Users", "₹ 1,499 / month", "Active", "M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z", "text-indigo-600"},
		{"WhatsApp Pack - 3000", "WhatsApp", "3,000 Messages", "₹ 899 / month", "Active", "M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z", "text-green-500"},
		{"Email Pack - 10000", "Email", "10,000 Emails", "₹ 599 / month", "Active", "M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z", "text-orange-500"},
		{"SMS Pack - 5000", "SMS", "5,000 SMS", "₹ 499 / month", "Active", "M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z", "text-green-500"},
		{"Storage 20 GB", "Storage", "20 GB", "₹ 799 / month", "Active", "M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z", "text-blue-500"},
		{"Storage 50 GB", "Storage", "50 GB", "₹ 1,499 / month", "Inactive", "M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z", "text-blue-500"},
		{"Voice Call Pack - 2000 Min", "Voice", "2,000 Minutes", "₹ 699 / month", "Active", "M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z", "text-orange-400"},
		{"Document Storage 50 GB", "Storage", "50 GB", "₹ 1,299 / month", "Inactive", "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z", "text-indigo-500"},
		{"Reports Pack", "Reports", "Advanced Reports", "₹ 999 / month", "Active", "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z", "text-indigo-500"},
		{"API Access Pack", "API", "API Access", "₹ 1,999 / month", "Active", "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4", "text-pink-500"},
	}

	for _, a := range addons {
		id := uuid.New()
		_, err := conn.Exec(context.Background(),
			"INSERT INTO mst_addon (addon_uuid, name, category, addon_limit, price, status, icon_svg, icon_color) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			id, a.Name, a.Category, a.Limit, a.Price, a.Status, a.IconSvg, a.IconColor)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Inserted mock addons")
}
