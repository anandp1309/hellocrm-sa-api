package handler

import (
	"net/http"

	"hellocrm-superadmin/server"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var e *echo.Echo
var dbPool *pgxpool.Pool
var initErr error

func init() {
	// Initialize the app once per lambda container
	e, dbPool, initErr = server.BuildApp()
}

// Handler is the serverless entrypoint for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	if initErr != nil {
		http.Error(w, "Startup Error: "+initErr.Error(), http.StatusInternalServerError)
		return
	}
	e.ServeHTTP(w, r)
}
