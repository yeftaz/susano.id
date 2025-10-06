package admin

import (
	"net/http"

	"github.com/yeftaz/susano.id/api/internal/service/admin"
	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

type UploadHandler struct {
	uploadService *admin.UploadService
	logger        *logger.Logger
}

func NewUploadHandler(uploadService *admin.UploadService, logger *logger.Logger) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		logger:        logger,
	}
}

// UploadAvatar handles POST /api/v1/admin/upload/avatar
func (h *UploadHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement avatar upload logic
	response.Success(w, nil, "Avatar upload endpoint - not implemented yet")
}

// DeleteAvatar handles DELETE /api/v1/admin/upload/avatar/{id}
func (h *UploadHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement avatar deletion logic
	response.Success(w, nil, "Avatar delete endpoint - not implemented yet")
}
