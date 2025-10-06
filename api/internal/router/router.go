package router

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/handler/shared"
	"github.com/yeftaz/susano.id/api/internal/middleware"
	"github.com/yeftaz/susano.id/api/pkg/logger"
)

// New creates and configures the main router
func New(cfg *config.Config, db *sql.DB, logger *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.RateLimiter(cfg))

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check
	healthHandler := shared.NewHealthHandler(db)
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Admin routes
	RegisterAdminRoutes(api, cfg, db, logger)

	// Store routes
	RegisterStoreRoutes(api, cfg, db, logger)

	return r
}

// ShowRoutes displays all registered routes (for make routes command)
func ShowRoutes(r *mux.Router) {
	println("╔════════╤═══════════════════════════════════════════════════╤═══════════════════════════════╤═══════════════════════╗")
	println("║ METHOD │ PATH                                              │ MIDDLEWARE                    │ HANDLER               ║")
	println("╠════════╪═══════════════════════════════════════════════════╪═══════════════════════════════╪═══════════════════════╣")

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		if len(methods) == 0 {
			return nil
		}

		method := methods[0]

		// Determine middleware based on path
		middleware := "-"
		if pathTemplate != "/api/v1/health" {
			middleware = "RateLimit"
			if pathTemplate != "/api/v1/admin/auth/login" && pathTemplate != "/api/v1/store/auth/login" && pathTemplate != "/api/v1/store/auth/register" {
				if len(pathTemplate) > 15 && pathTemplate[:15] == "/api/v1/admin/" {
					middleware = "RateLimit, AdminAuth"
				} else if len(pathTemplate) > 15 && pathTemplate[:15] == "/api/v1/store/" {
					middleware = "RateLimit, CustomerAuth"
				}
			}
		}

		// Determine handler name from path
		handler := getHandlerName(pathTemplate)

		printf("║ %-6s │ %-49s │ %-29s │ %-21s ║\n", method, pathTemplate, middleware, handler)
		return nil
	})

	println("╚════════╧═══════════════════════════════════════════════════╧═══════════════════════════════╧═══════════════════════╝")
}

func getHandlerName(path string) string {
	handlers := map[string]string{
		"/api/v1/health":                   "HealthCheck",
		"/api/v1/admin/auth/login":         "Login",
		"/api/v1/admin/auth/logout":        "Logout",
		"/api/v1/admin/auth/me":            "GetCurrentUser",
		"/api/v1/admin/auth/refresh":       "RefreshSession",
		"/api/v1/admin/admins":             "GetAll/Create",
		"/api/v1/admin/admins/{id}":        "GetByID/Update/Delete",
		"/api/v1/admin/dashboard/stats":    "GetStats",
		"/api/v1/admin/upload/avatar":      "UploadAvatar",
		"/api/v1/admin/upload/avatar/{id}": "DeleteAvatar",
		"/api/v1/store/auth/login":         "Login",
		"/api/v1/store/auth/register":      "Register",
		"/api/v1/store/auth/logout":        "Logout",
		"/api/v1/store/profile":            "GetProfile/UpdateProfile",
	}

	if handler, ok := handlers[path]; ok {
		return handler
	}
	return "Unknown"
}

func printf(format string, args ...interface{}) {
	print(sprintf(format, args...))
}

func sprintf(format string, args ...interface{}) string {
	// Simple sprintf implementation for the table
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			result = replaceFirst(result, "%s", v)
			result = replaceFirst(result, "%-6s", padRight(v, 6))
			result = replaceFirst(result, "%-49s", padRight(v, 49))
			result = replaceFirst(result, "%-29s", padRight(v, 29))
			result = replaceFirst(result, "%-21s", padRight(v, 21))
		}
	}
	return result
}

func replaceFirst(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + string(make([]byte, length-len(s)))
}
