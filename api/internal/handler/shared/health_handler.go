package shared

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/yeftaz/susano.id/api/pkg/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

type HealthResponse struct {
	Status   string         `json:"status"`
	Service  string         `json:"service"`
	Database DatabaseHealth `json:"database"`
	Uptime   string         `json:"uptime,omitempty"`
}

type DatabaseHealth struct {
	Status       string `json:"status"`
	ResponseTime string `json:"response_time,omitempty"`
}

// HealthCheck handles GET /api/v1/health
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check database connection
	dbStatus := "ok"
	var dbResponseTime time.Duration
	start := time.Now()

	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "error"
		response.Error(w, http.StatusServiceUnavailable, "Service unhealthy")
		return
	}

	dbResponseTime = time.Since(start)

	healthData := HealthResponse{
		Status:  "ok",
		Service: "susano-api",
		Database: DatabaseHealth{
			Status:       dbStatus,
			ResponseTime: dbResponseTime.String(),
		},
	}

	response.Success(w, healthData, "Service is healthy")
}
