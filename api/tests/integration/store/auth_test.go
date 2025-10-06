package store_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/database"
	"github.com/yeftaz/susano.id/api/internal/router"
	"github.com/yeftaz/susano.id/api/pkg/logger"
)

func setupTestRouter(t *testing.T) *http.Handler {
	cfg := &config.Config{
		AppEnv:          "test",
		DBHost:          "localhost",
		DBPort:          "5432",
		DBName:          "susano_test",
		DBUser:          "root",
		DBPassword:      "",
		DBSSLMode:       "disable",
		SessionLifetime: 720 * 3600,
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	testLogger := logger.New(cfg)
	r := router.New(cfg, db, testLogger)
	var handler http.Handler = r
	return &handler
}

func TestCustomerRegister(t *testing.T) {
	handler := setupTestRouter(t)

	t.Run("Successful Registration", func(t *testing.T) {
		registerData := map[string]string{
			"email":    "customer@example.com",
			"password": "password123",
			"name":     "Test Customer",
		}

		body, _ := json.Marshal(registerData)
		req := httptest.NewRequest("POST", "/api/v1/store/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		// TODO: Implement when customer registration is ready
		t.Skip("Customer registration not yet implemented")
	})
}

func TestCustomerLogin(t *testing.T) {
	handler := setupTestRouter(t)

	t.Run("Customer Login", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "customer@example.com",
			"password": "password123",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/store/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		// TODO: Implement when customer login is ready
		t.Skip("Customer login not yet implemented")
	})
}
