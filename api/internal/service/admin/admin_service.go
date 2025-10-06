package admin

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/yeftaz/susano.id/api/internal/domain/admin"
	adminRepo "github.com/yeftaz/susano.id/api/internal/repository/admin"
)

type AdminService struct {
	adminRepo *adminRepo.AdminRepository
}

func NewAdminService(adminRepo *adminRepo.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

// GetAll retrieves all admins with pagination and filtering
func (s *AdminService) GetAll(ctx context.Context, page, limit int, search, role string) ([]*admin.Admin, int, error) {
	admins, total, err := s.adminRepo.GetAll(ctx, page, limit, search, role)
	if err != nil {
		return nil, 0, err
	}

	// Clear passwords before returning
	for _, a := range admins {
		a.Password = ""
	}

	return admins, total, nil
}

// GetByID retrieves an admin by ID
func (s *AdminService) GetByID(ctx context.Context, id string) (*admin.Admin, error) {
	admin, err := s.adminRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	admin.Password = ""

	return admin, nil
}

// Create creates a new admin
func (s *AdminService) Create(ctx context.Context, email, password, name, role string) (*admin.Admin, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create admin
	admin, err := s.adminRepo.Create(ctx, email, string(hashedPassword), name, role)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	admin.Password = ""

	return admin, nil
}

// Update updates an admin
func (s *AdminService) Update(ctx context.Context, id string, email, name, role string) (*admin.Admin, error) {
	var emailPtr, namePtr, rolePtr *string

	if email != "" {
		emailPtr = &email
	}
	if name != "" {
		namePtr = &name
	}
	if role != "" {
		rolePtr = &role
	}

	admin, err := s.adminRepo.Update(ctx, id, emailPtr, namePtr, rolePtr)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	admin.Password = ""

	return admin, nil
}

// Delete soft deletes an admin
func (s *AdminService) Delete(ctx context.Context, id string) error {
	return s.adminRepo.Delete(ctx, id)
}

// UpdateAvatar updates admin avatar
func (s *AdminService) UpdateAvatar(ctx context.Context, id, avatarPath string) error {
	return s.adminRepo.UpdateAvatarPath(ctx, id, avatarPath)
}
