package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"hellocrm-superadmin/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Force load the exact absolute path to the API project's .env file (for local dev)
	_ = godotenv.Overload("C:/hellocrm-superadmin/EZ_Engineering_OS_v1.0_Final/.env")

	e, dbPool := server.BuildApp()
	defer dbPool.Close()

	// Start server gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	log.Println("shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
