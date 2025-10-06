package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminCRUD(t *testing.T) {
	handler := setupTestRouter(t)

	// Get session token first
	sessionCookie := loginAndGetSessionCookie(t, handler)

	// Test Get All Admins
	t.Run("Get All Admins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/admins", nil)
		req.AddCookie(sessionCookie)

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	// Test Create Admin
	t.Run("Create Admin", func(t *testing.T) {
		adminData := map[string]string{
			"email":    "newadmin@susano.id",
			"password": "password123",
			"name":     "New Admin",
			"role":     "admin",
		}

		body, _ := json.Marshal(adminData)
		req := httptest.NewRequest("POST", "/api/v1/admin/admins", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		rr := httptest.NewRecorder()
		(*handler).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}

func loginAndGetSessionCookie(t *testing.T, handler *http.Handler) *http.Cookie {
	loginData := map[string]string{
		"email":    "admin@susano.id",
		"password": "admin1234",
	}

	body, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	(*handler).ServeHTTP(rr, req)

	cookies := rr.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			return cookie
		}
	}

	t.Fatal("No session cookie found")
	return nil
}
