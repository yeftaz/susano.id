package middleware

import (
	"context"
	"net/http"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/service/store"
	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

// Context keys for customer
type customerContextKey string

const (
	contextKeyCustomer   customerContextKey = "customer"
	contextKeyCustomerID customerContextKey = "customer_id"
)

// CustomerAuth middleware verifies customer session
func CustomerAuth(authService *store.AuthService, cfg *config.Config, logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session token from cookie
			cookie, err := r.Cookie("customer_session_token")
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized: No session found")
				return
			}

			// Verify session and get customer
			customer, err := authService.VerifySession(r.Context(), cookie.Value, cfg.SessionLifetime)
			if err != nil {
				logger.Error("Customer session verification failed", "error", err)
				response.Error(w, http.StatusUnauthorized, "Unauthorized: Invalid or expired session")
				return
			}

			// Check if customer is active and not deleted
			if !customer.CanPurchase() {
				response.Error(w, http.StatusForbidden, "Account is inactive or deleted")
				return
			}

			// Add customer to request context
			ctx := context.WithValue(r.Context(), contextKeyCustomer, customer)
			ctx = context.WithValue(ctx, contextKeyCustomerID, customer.ID.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
