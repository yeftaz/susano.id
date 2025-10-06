package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yeftaz/susano.id/api/internal/config"
	adminHandler "github.com/yeftaz/susano.id/api/internal/handler/admin"
	"github.com/yeftaz/susano.id/api/internal/middleware"
	adminRepo "github.com/yeftaz/susano.id/api/internal/repository/admin"
	adminService "github.com/yeftaz/susano.id/api/internal/service/admin"
	"github.com/yeftaz/susano.id/api/pkg/logger"
)

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *mux.Router, cfg *config.Config, db *sql.DB, logger *logger.Logger) {
	// Initialize repositories
	adminRepository := adminRepo.NewAdminRepository(db)
	sessionRepository := adminRepo.NewSessionRepository(db)

	// Initialize services
	authService := adminService.NewAuthService(adminRepository, sessionRepository)
	adminSvc := adminService.NewAdminService(adminRepository)
	uploadService := adminService.NewUploadService()

	// Initialize handlers
	authHandler := adminHandler.NewAuthHandler(authService, logger, cfg)
	adminHdlr := adminHandler.NewAdminHandler(adminSvc, logger)
	dashboardHandler := adminHandler.NewDashboardHandler(logger)
	uploadHandler := adminHandler.NewUploadHandler(uploadService, logger)

	// Auth middleware
	adminAuth := middleware.AdminAuth(authService, cfg, logger)

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()

	// Auth routes (public)
	admin.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Auth routes (protected)
	admin.Handle("/auth/logout", adminAuth(http.HandlerFunc(authHandler.Logout))).Methods("POST")
	admin.Handle("/auth/me", adminAuth(http.HandlerFunc(authHandler.GetCurrentUser))).Methods("GET")
	admin.Handle("/auth/refresh", adminAuth(http.HandlerFunc(authHandler.RefreshSession))).Methods("POST")

	// Admin CRUD routes (protected)
	admin.Handle("/admins", adminAuth(http.HandlerFunc(adminHdlr.GetAll))).Methods("GET")
	admin.Handle("/admins", adminAuth(http.HandlerFunc(adminHdlr.Create))).Methods("POST")
	admin.Handle("/admins/{id}", adminAuth(http.HandlerFunc(adminHdlr.GetByID))).Methods("GET")
	admin.Handle("/admins/{id}", adminAuth(http.HandlerFunc(adminHdlr.Update))).Methods("PATCH")
	admin.Handle("/admins/{id}", adminAuth(http.HandlerFunc(adminHdlr.Delete))).Methods("DELETE")

	// Dashboard routes (protected)
	admin.Handle("/dashboard/stats", adminAuth(http.HandlerFunc(dashboardHandler.GetStats))).Methods("GET")

	// Upload routes (protected)
	admin.Handle("/upload/avatar", adminAuth(http.HandlerFunc(uploadHandler.UploadAvatar))).Methods("POST")
	admin.Handle("/upload/avatar/{id}", adminAuth(http.HandlerFunc(uploadHandler.DeleteAvatar))).Methods("DELETE")
}
