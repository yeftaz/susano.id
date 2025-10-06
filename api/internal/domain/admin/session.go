package admin

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an admin session entity
type Session struct {
	ID             uuid.UUID `json:"id"`
	AdminID        uuid.UUID `json:"admin_id"`
	Token          string    `json:"-"` // Never expose token in JSON
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	LastActivityAt time.Time `json:"last_activity_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// IsExpired checks if the session has expired based on session lifetime
func (s *Session) IsExpired(sessionLifetime time.Duration) bool {
	expiresAt := s.CreatedAt.Add(sessionLifetime)
	return time.Now().After(expiresAt)
}

// ShouldRefresh checks if session should be refreshed (last activity > 15 minutes ago)
func (s *Session) ShouldRefresh() bool {
	return time.Since(s.LastActivityAt) > 15*time.Minute
}
