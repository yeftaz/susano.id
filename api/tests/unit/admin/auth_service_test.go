package admin_test

import (
	"context"
	"testing"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/database"
	adminRepo "github.com/yeftaz/susano.id/api/internal/repository/admin"
	adminService "github.com/yeftaz/susano.id/api/internal/service/admin"
)

func setupAuthService(t *testing.T) *adminService.AuthService {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "susano_test",
		DBUser:     "root",
		DBPassword: "",
		DBSSLMode:  "disable",
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	adminRepository := adminRepo.NewAdminRepository(db)
	sessionRepository := adminRepo.NewSessionRepository(db)

	return adminService.NewAuthService(adminRepository, sessionRepository)
}

func TestLogin(t *testing.T) {
	authService := setupAuthService(t)
	ctx := context.Background()

	t.Run("Successful Login", func(t *testing.T) {
		admin, session, err := authService.Login(ctx, "admin@susano.id", "admin1234", "127.0.0.1", "test-agent")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if admin == nil {
			t.Error("Expected admin to be returned")
		}

		if session == nil {
			t.Error("Expected session to be returned")
		}

		if admin.Email != "admin@susano.id" {
			t.Errorf("Expected email admin@susano.id, got %s", admin.Email)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, _, err := authService.Login(ctx, "admin@susano.id", "wrongpassword", "127.0.0.1", "test-agent")

		if err == nil {
			t.Error("Expected error for invalid password")
		}
	})

	t.Run("Invalid Email", func(t *testing.T) {
		_, _, err := authService.Login(ctx, "nonexistent@susano.id", "admin1234", "127.0.0.1", "test-agent")

		if err == nil {
			t.Error("Expected error for invalid email")
		}
	})
}
