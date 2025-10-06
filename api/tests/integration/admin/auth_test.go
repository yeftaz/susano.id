package admin_test

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
	// Load test config
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

	// Connect to test database
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Initialize logger
	testLogger := logger.New(cfg)

	// Create router
	r := router.New(cfg, db, testLogger)
	var handler http.Handler = r
	return &handler
}

func TestAdminLogin(t *testing.T) {
	handler := setupTestRouter(t)

	// Test successful login
	t.Run("Successful Login", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "admin@susano.id",
			"password": "admin1234",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check response body
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		if success, ok := response["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", response["success"])
		}
	})

	// Test invalid credentials
	t.Run("Invalid Credentials", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "admin@susano.id",
			"password": "wrongpassword",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	// Test validation error
	t.Run("Validation Error", func(t *testing.T) {
		loginData := map[string]string{
			"email": "invalid-email",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
		}
	})
}

func TestAdminLogout(t *testing.T) {
	handler := setupTestRouter(t)

	// First, login to get session token
	loginData := map[string]string{
		"email":    "admin@susano.id",
		"password": "admin1234",
	}

	body, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	(*handler).ServeHTTP(rr, req)

	// Extract session cookie
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Fatal("No session cookie found after login")
	}

	// Test logout
	t.Run("Successful Logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/auth/logout", nil)
		req.AddCookie(sessionCookie)

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
