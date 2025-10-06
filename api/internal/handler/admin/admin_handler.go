package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yeftaz/susano.id/api/internal/service/admin"
	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
	"github.com/yeftaz/susano.id/api/pkg/validator"
)

type AdminHandler struct {
	adminService *admin.AdminService
	logger       *logger.Logger
}

func NewAdminHandler(adminService *admin.AdminService, logger *logger.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       logger,
	}
}

type CreateAdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin cashier"`
}

type UpdateAdminRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty"`
	Role  string `json:"role" validate:"omitempty,oneof=super_admin admin cashier"`
}

// GetAll handles GET /api/v1/admin/admins
func (h *AdminHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	search := r.URL.Query().Get("search")
	role := r.URL.Query().Get("role")

	// Get admins
	admins, total, err := h.adminService.GetAll(r.Context(), page, limit, search, role)
	if err != nil {
		h.logger.Error("Failed to get admins", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve admins")
		return
	}

	response.SuccessWithMeta(w, admins, "Admins retrieved successfully", map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetByID handles GET /api/v1/admin/admins/{id}
func (h *AdminHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	admin, err := h.adminService.GetByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusNotFound, "Admin not found")
			return
		}
		h.logger.Error("Failed to get admin", "id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve admin")
		return
	}

	response.Success(w, admin, "Admin retrieved successfully")
}

// Create handles POST /api/v1/admin/admins
func (h *AdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Create admin
	admin, err := h.adminService.Create(r.Context(), req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		h.logger.Error("Failed to create admin", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to create admin")
		return
	}

	h.logger.Info("Admin created", "admin_id", admin.ID, "email", admin.Email)
	response.Success(w, admin, "Admin created successfully")
}

// Update handles PATCH /api/v1/admin/admins/{id}
func (h *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Update admin
	admin, err := h.adminService.Update(r.Context(), id, req.Email, req.Name, req.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusNotFound, "Admin not found")
			return
		}
		h.logger.Error("Failed to update admin", "id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to update admin")
		return
	}

	h.logger.Info("Admin updated", "admin_id", id)
	response.Success(w, admin, "Admin updated successfully")
}

// Delete handles DELETE /api/v1/admin/admins/{id}
func (h *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Soft delete admin
	if err := h.adminService.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusNotFound, "Admin not found")
			return
		}
		h.logger.Error("Failed to delete admin", "id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to delete admin")
		return
	}

	h.logger.Info("Admin deleted", "admin_id", id)
	response.Success(w, nil, "Admin deleted successfully")
}
