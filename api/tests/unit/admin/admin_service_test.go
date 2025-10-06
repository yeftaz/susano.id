package admin_test

import (
	"context"
	"testing"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/database"
	adminRepo "github.com/yeftaz/susano.id/api/internal/repository/admin"
	adminService "github.com/yeftaz/susano.id/api/internal/service/admin"
)

func setupAdminService(t *testing.T) *adminService.AdminService {
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
	return adminService.NewAdminService(adminRepository)
}

func TestGetAllAdmins(t *testing.T) {
	service := setupAdminService(t)
	ctx := context.Background()

	t.Run("Get All Admins", func(t *testing.T) {
		admins, total, err := service.GetAll(ctx, 1, 10, "", "")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if admins == nil {
			t.Error("Expected admins to be returned")
		}

		if total == 0 {
			t.Error("Expected at least one admin in database")
		}
	})
}

func TestCreateAdmin(t *testing.T) {
	service := setupAdminService(t)
	ctx := context.Background()

	t.Run("Create Admin", func(t *testing.T) {
		admin, err := service.Create(ctx, "test@susano.id", "password123", "Test Admin", "admin")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if admin == nil {
			t.Error("Expected admin to be returned")
		}

		if admin.Email != "test@susano.id" {
			t.Errorf("Expected email test@susano.id, got %s", admin.Email)
		}

		// Cleanup: delete test admin
		_ = service.Delete(ctx, admin.ID.String())
	})
}
