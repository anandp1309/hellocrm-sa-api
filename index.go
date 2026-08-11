package api

import (
	"net/http"

	"hellocrm-superadmin/server"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var e *echo.Echo
var dbPool *pgxpool.Pool

func init() {
	// Initialize the app once per lambda container
	e, dbPool = server.BuildApp()
}

// Handler is the serverless entrypoint for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}
