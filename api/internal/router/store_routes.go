package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yeftaz/susano.id/api/internal/config"
	storeHandler "github.com/yeftaz/susano.id/api/internal/handler/store"
	"github.com/yeftaz/susano.id/api/internal/middleware"
	storeRepo "github.com/yeftaz/susano.id/api/internal/repository/store"
	storeService "github.com/yeftaz/susano.id/api/internal/service/store"
	"github.com/yeftaz/susano.id/api/pkg/logger"
)

// RegisterStoreRoutes registers all store (customer) routes
func RegisterStoreRoutes(r *mux.Router, cfg *config.Config, db *sql.DB, logger *logger.Logger) {
	// Initialize repositories
	customerRepository := storeRepo.NewCustomerRepository(db)
	sessionRepository := storeRepo.NewSessionRepository(db)

	// Initialize services
	authService := storeService.NewAuthService(customerRepository, sessionRepository)
	customerService := storeService.NewCustomerService(customerRepository)

	// Initialize handlers
	authHandler := storeHandler.NewAuthHandler(logger)
	customerHandler := storeHandler.NewCustomerHandler(logger)

	// Auth middleware
	customerAuth := middleware.CustomerAuth(authService, cfg, logger)

	// Store routes
	store := r.PathPrefix("/store").Subrouter()

	// Auth routes (public)
	store.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	store.HandleFunc("/auth/register", authHandler.Register).Methods("POST")

	// Auth routes (protected)
	store.Handle("/auth/logout", customerAuth(http.HandlerFunc(authHandler.Logout))).Methods("POST")

	// Customer profile routes (protected)
	store.Handle("/profile", customerAuth(http.HandlerFunc(customerHandler.GetProfile))).Methods("GET")
	store.Handle("/profile", customerAuth(http.HandlerFunc(customerHandler.UpdateProfile))).Methods("PATCH")

	// Suppress unused variable warnings for now
	_ = authService
	_ = customerService
}
