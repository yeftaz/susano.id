package admin

import (
	"encoding/json"
	"net/http"

	"github.com/yeftaz/susano.id/api/internal/config"
	adminDomain "github.com/yeftaz/susano.id/api/internal/domain/admin"
	"github.com/yeftaz/susano.id/api/internal/service/admin"
	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
	"github.com/yeftaz/susano.id/api/pkg/validator"
)

type AuthHandler struct {
	authService *admin.AuthService
	logger      *logger.Logger
	config      *config.Config
}

func NewAuthHandler(authService *admin.AuthService, logger *logger.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		config:      cfg,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Admin interface{} `json:"admin"`
	Token string      `json:"token"`
}

// Login handles POST /api/v1/admin/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Authenticate admin
	admin, session, err := h.authService.Login(r.Context(), req.Email, req.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		h.logger.Error("Login failed", "email", req.Email, "error", err)
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	h.logger.Info("Admin logged in", "admin_id", admin.ID, "email", admin.Email)

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.SessionSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.config.SessionLifetime.Seconds()),
		Domain:   h.config.SessionDomain,
	})

	response.Success(w, LoginResponse{
		Admin: admin,
		Token: session.Token,
	}, "Login successful")
}

// Logout handles POST /api/v1/admin/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "No session found")
		return
	}

	// Delete session
	if err := h.authService.Logout(r.Context(), cookie.Value); err != nil {
		h.logger.Error("Logout failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.SessionSecure,
		MaxAge:   -1,
		Domain:   h.config.SessionDomain,
	})

	response.Success(w, nil, "Logout successful")
}

// GetCurrentUser handles GET /api/v1/admin/auth/me
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get admin from context (set by auth middleware)
	admin, ok := r.Context().Value("admin").(*adminDomain.Admin)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response.Success(w, admin, "Admin retrieved successfully")
}

// RefreshSession handles POST /api/v1/admin/auth/refresh
func (h *AuthHandler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "No session found")
		return
	}

	// Refresh session
	if err := h.authService.RefreshSession(r.Context(), cookie.Value); err != nil {
		h.logger.Error("Session refresh failed", "error", err)
		response.Error(w, http.StatusUnauthorized, "Failed to refresh session")
		return
	}

	response.Success(w, nil, "Session refreshed successfully")
}
