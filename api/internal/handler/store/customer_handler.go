package store

import (
	"net/http"

	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

type CustomerHandler struct {
	logger *logger.Logger
}

func NewCustomerHandler(logger *logger.Logger) *CustomerHandler {
	return &CustomerHandler{
		logger: logger,
	}
}

// GetProfile handles GET /api/v1/store/profile
func (h *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get customer profile
	response.Success(w, nil, "Customer profile endpoint - not implemented yet")
}

// UpdateProfile handles PATCH /api/v1/store/profile
func (h *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update customer profile
	response.Success(w, nil, "Update customer profile endpoint - not implemented yet")
}
