package store

import (
	"net/http"

	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

type AuthHandler struct {
	logger *logger.Logger
}

func NewAuthHandler(logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		logger: logger,
	}
}

// Login handles POST /api/v1/store/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement customer login
	response.Success(w, nil, "Customer login endpoint - not implemented yet")
}

// Register handles POST /api/v1/store/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement customer registration
	response.Success(w, nil, "Customer register endpoint - not implemented yet")
}

// Logout handles POST /api/v1/store/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement customer logout
	response.Success(w, nil, "Customer logout endpoint - not implemented yet")
}
