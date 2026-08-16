package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"hellocrm-superadmin/internal/auth"
	"hellocrm-superadmin/internal/customer"
	"hellocrm-superadmin/internal/identity"
	"hellocrm-superadmin/internal/mastertype"
	"hellocrm-superadmin/internal/gateway"
	"hellocrm-superadmin/internal/invoice"
	"hellocrm-superadmin/internal/admin"
	"hellocrm-superadmin/internal/role"
	"hellocrm-superadmin/internal/auditlog"
	"hellocrm-superadmin/internal/jobs"
	"hellocrm-superadmin/internal/cron"
	"hellocrm-superadmin/internal/health"
	"hellocrm-superadmin/internal/templates"
	"hellocrm-superadmin/internal/announcements"
	"hellocrm-superadmin/internal/errorlog"
	"hellocrm-superadmin/internal/addon"
	"hellocrm-superadmin/internal/payment"
	"hellocrm-superadmin/internal/plan"
	"hellocrm-superadmin/internal/platform/database/db"
	"hellocrm-superadmin/internal/platform/websocket"
	"hellocrm-superadmin/internal/subscription"
	"hellocrm-superadmin/internal/tenant"
	"hellocrm-superadmin/internal/universal"
	"hellocrm-superadmin/internal/usage"
	"hellocrm-superadmin/internal/dashboard"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BuildApp() (*echo.Echo, *pgxpool.Pool) {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS Setup
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowOrigins := []string{"*"} // Default
	if corsOrigins != "" {
		allowOrigins = strings.Split(corsOrigins, ",")
		for i, origin := range allowOrigins {
			allowOrigins[i] = strings.TrimSpace(origin)
		}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Requested-With", "Cache-Control"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// Infrastructure
	wsHub := websocket.NewHub()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Provide a fallback or log fatal if required
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Dependency Injection
	queries := db.New(dbPool)

	// Auth & RBAC
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_key_change_me"
	}
	
	enforcer, err := auth.NewCasbinEnforcer(queries)
	if err != nil {
		log.Fatalf("Unable to initialize Casbin enforcer: %v", err)
	}

	authService := auth.NewService(queries, []byte(jwtSecret))
	authMW := auth.NewMiddleware(authService, enforcer)

	roleRepo := role.NewRepository(queries)
	roleQuery := role.NewQueryService(roleRepo)
	roleCommand := role.NewCommandService(roleRepo)
	roleHandler := role.NewHandler(roleQuery, roleCommand)

	authHandler := auth.NewHandler(authService, roleQuery)

	// Features
	tenantRepo := tenant.NewRepository(queries)
	tenantCmd := tenant.NewCommandService(tenantRepo)
	tenantQuery := tenant.NewQueryService(tenantRepo)
	tenantHandler := tenant.NewHandler(tenantCmd, tenantQuery)

	identityRepo := identity.NewRepository(queries)
	identityCmd := identity.NewUserCommandService(identityRepo)
	identityQuery := identity.NewUserQueryService(identityRepo)
	identityHandler := identity.NewHandler(identityCmd, identityQuery)

	customerRepo := customer.NewRepository(queries)
	customerCmd := customer.NewCommandService(customerRepo)
	customerQuery := customer.NewQueryService(customerRepo)
	customerHandler := customer.NewHandler(customerCmd, customerQuery)

	subscriptionRepo := subscription.NewRepository(queries)
	subscriptionCmd := subscription.NewCommandService(subscriptionRepo)
	subscriptionQuery := subscription.NewQueryService(subscriptionRepo)
	subscriptionHandler := subscription.NewHandler(subscriptionCmd, subscriptionQuery)

	planRepo := plan.NewRepository(queries)
	planCmd := plan.NewCommandService(planRepo)
	planQuery := plan.NewQueryService(planRepo)
	planHandler := plan.NewHandler(planCmd, planQuery)

	masterTypeRepo := mastertype.NewRepository(queries)
	masterTypeService := mastertype.NewService(masterTypeRepo)
	masterTypeHandler := mastertype.NewHandler(masterTypeService)

	universalRepo := universal.NewRepository(queries)
	universalService := universal.NewService(universalRepo)
	universalHandler := universal.NewHandler(universalService)

	paymentRepo := payment.NewRepository(queries)
	paymentQuery := payment.NewQueryService(paymentRepo)
	paymentHandler := payment.NewHandler(paymentQuery)

	gatewayRepo := gateway.NewRepository(queries)
	gatewayQuery := gateway.NewQueryService(gatewayRepo)
	gatewayHandler := gateway.NewHandler(gatewayQuery)

	invoiceRepo := invoice.NewRepository(queries)
	invoiceQuery := invoice.NewQueryService(invoiceRepo)
	invoiceHandler := invoice.NewHandler(invoiceQuery)

	adminRepo := admin.NewRepository(queries)
	adminQuery := admin.NewQueryService(adminRepo)
	adminCommand := admin.NewCommandService(adminRepo)
	adminHandler := admin.NewHandler(adminQuery, adminCommand)


	auditLogHandler := auditlog.NewHandler()
	jobsHandler := jobs.NewHandler()
	cronHandler := cron.NewHandler()
	sysHealthHandler := health.NewHandler()
	templatesHandler := templates.NewHandler()
	announcementsHandler := announcements.NewHandler()
	errorLogHandler := errorlog.NewHandler()

	// Routes
	v1 := e.Group("/api/v1")
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	tenantHandler.RegisterRoutes(v1)
	identityHandler.RegisterRoutes(v1, authMW)
	authHandler.RegisterRoutes(v1, authMW)
	customerHandler.RegisterRoutes(v1)
	subscriptionHandler.RegisterRoutes(v1)
	planHandler.RegisterRoutes(v1)
	masterTypeHandler.RegisterRoutes(v1)
	universalHandler.RegisterRoutes(v1)
	paymentHandler.RegisterRoutes(v1)
	gatewayHandler.RegisterRoutes(v1)
	invoiceHandler.RegisterRoutes(v1)
	adminHandler.RegisterRoutes(v1)
	roleHandler.RegisterRoutes(v1)
	
	auditLogHandler.RegisterRoutes(v1)
	jobsHandler.RegisterRoutes(v1)
	cronHandler.RegisterRoutes(v1)
	sysHealthHandler.RegisterRoutes(v1)
	templatesHandler.RegisterRoutes(v1)
	announcementsHandler.RegisterRoutes(v1)
	errorLogHandler.RegisterRoutes(v1)
	
	// WebSocket Route
	e.GET("/ws", wsHub.HandleWebSocket)

	usageQuery := usage.NewQueryService(queries)
	usageHandler := usage.NewHandler(usageQuery)
	usageHandler.RegisterRoutes(v1)

	dashboardQuery := dashboard.NewQueryService(queries)
	dashboardHandler := dashboard.NewHandler(dashboardQuery)
	dashboardHandler.RegisterRoutes(v1)

	addonRepo := addon.NewRepository(queries)
	addonCmd := addon.NewCommandService(addonRepo)
	addonQuery := addon.NewQueryService(addonRepo)
	addonHandler := addon.NewHandler(addonCmd, addonQuery)
	addonHandler.RegisterRoutes(v1)

	return e, dbPool
}
