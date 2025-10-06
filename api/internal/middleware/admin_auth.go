package middleware

import (
	"context"
	"net/http"

	"github.com/yeftaz/susano.id/api/internal/config"
	adminDomain "github.com/yeftaz/susano.id/api/internal/domain/admin"
	"github.com/yeftaz/susano.id/api/internal/service/admin"
	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

// Context keys for admin
type contextKey string

const (
	contextKeyAdmin   contextKey = "admin"
	contextKeyAdminID contextKey = "admin_id"
)

// AdminAuth middleware verifies admin session
func AdminAuth(authService *admin.AuthService, cfg *config.Config, logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session token from cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized: No session found")
				return
			}

			// Verify session and get admin
			adminUser, err := authService.VerifySession(r.Context(), cookie.Value, cfg.SessionLifetime)
			if err != nil {
				logger.Error("Session verification failed", "error", err)
				response.Error(w, http.StatusUnauthorized, "Unauthorized: Invalid or expired session")
				return
			}

			// Check if admin is active and not deleted
			if !adminUser.CanAccessAdminPanel() {
				response.Error(w, http.StatusForbidden, "Account is inactive or deleted")
				return
			}

			// Add admin to request context
			ctx := context.WithValue(r.Context(), contextKeyAdmin, adminUser)
			ctx = context.WithValue(ctx, contextKeyAdminID, adminUser.ID.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware checks if admin has required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			adminUser, ok := r.Context().Value(contextKeyAdmin).(*adminDomain.Admin)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Check if admin has any of the required roles
			hasRole := false
			for _, role := range roles {
				if string(adminUser.Role) == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Error(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
