package store

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer user entity
type Customer struct {
	ID                     uuid.UUID  `json:"id"`
	Email                  string     `json:"email"`
	Password               string     `json:"-"` // Never expose password in JSON
	Name                   string     `json:"name"`
	AvatarPath             *string    `json:"avatar_path,omitempty"`
	IsActive               bool       `json:"is_active"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at,omitempty"`
	TwoFactorSecret        *string    `json:"-"` // Never expose 2FA secret
	TwoFactorRecoveryCodes *string    `json:"-"` // Never expose recovery codes
	TwoFactorConfirmedAt   *time.Time `json:"two_factor_confirmed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

// IsDeleted checks if the customer is soft deleted
func (c *Customer) IsDeleted() bool {
	return c.DeletedAt != nil
}

// HasTwoFactor checks if customer has 2FA enabled
func (c *Customer) HasTwoFactor() bool {
	return c.TwoFactorConfirmedAt != nil
}

// CanPurchase checks if customer can make purchases
func (c *Customer) CanPurchase() bool {
	return c.IsActive && !c.IsDeleted()
}

// IsEmailVerified checks if customer email is verified
func (c *Customer) IsEmailVerified() bool {
	return c.EmailVerifiedAt != nil
}
