package admin

import (
	"net/http"

	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

type DashboardHandler struct {
	logger *logger.Logger
}

func NewDashboardHandler(logger *logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		logger: logger,
	}
}

// GetStats handles GET /api/v1/admin/dashboard/stats
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual stats logic
	stats := map[string]interface{}{
		"total_admins":    3,
		"total_customers": 0,
		"total_orders":    0,
		"total_revenue":   0,
	}

	response.Success(w, stats, "Dashboard stats retrieved successfully")
}
