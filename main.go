package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"hellocrm-superadmin/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file in the current directory
	_ = godotenv.Overload(".env")

	e, dbPool, err := server.BuildApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
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
