package admin

import (
	"time"

	"github.com/google/uuid"
)

// AdminRole represents the role of an admin user
type AdminRole string

const (
	RoleSuperAdmin AdminRole = "super_admin"
	RoleAdmin      AdminRole = "admin"
	RoleCashier    AdminRole = "cashier"
)

// Admin represents an admin user entity
type Admin struct {
	ID                     uuid.UUID  `json:"id"`
	Email                  string     `json:"email"`
	Password               string     `json:"-"` // Never expose password in JSON
	Name                   string     `json:"name"`
	AvatarPath             *string    `json:"avatar_path,omitempty"`
	Role                   AdminRole  `json:"role"`
	IsActive               bool       `json:"is_active"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at,omitempty"`
	TwoFactorSecret        *string    `json:"-"` // Never expose 2FA secret
	TwoFactorRecoveryCodes *string    `json:"-"` // Never expose recovery codes
	TwoFactorConfirmedAt   *time.Time `json:"two_factor_confirmed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

// IsDeleted checks if the admin is soft deleted
func (a *Admin) IsDeleted() bool {
	return a.DeletedAt != nil
}

// HasTwoFactor checks if admin has 2FA enabled
func (a *Admin) HasTwoFactor() bool {
	return a.TwoFactorConfirmedAt != nil
}

// CanAccessAdminPanel checks if admin can access the admin panel
func (a *Admin) CanAccessAdminPanel() bool {
	return a.IsActive && !a.IsDeleted()
}

// IsSuperAdmin checks if the admin has super admin role
func (a *Admin) IsSuperAdmin() bool {
	return a.Role == RoleSuperAdmin
}

// IsAdmin checks if the admin has admin role
func (a *Admin) IsAdmin() bool {
	return a.Role == RoleAdmin
}

// IsCashier checks if the admin has cashier role
func (a *Admin) IsCashier() bool {
	return a.Role == RoleCashier
}
