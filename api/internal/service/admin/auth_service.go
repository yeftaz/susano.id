package admin

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/yeftaz/susano.id/api/internal/domain/admin"
	adminRepo "github.com/yeftaz/susano.id/api/internal/repository/admin"
)

type AuthService struct {
	adminRepo   *adminRepo.AdminRepository
	sessionRepo *adminRepo.SessionRepository
}

func NewAuthService(adminRepo *adminRepo.AdminRepository, sessionRepo *adminRepo.SessionRepository) *AuthService {
	return &AuthService{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
	}
}

// Login authenticates an admin and creates a session
func (s *AuthService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*admin.Admin, *admin.Session, error) {
	// Find admin by email
	admin, err := s.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, nil, err
	}

	// Check if admin can access admin panel
	if !admin.CanAccessAdminPanel() {
		return nil, nil, err
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, nil, err
	}

	// Create session
	session, err := s.sessionRepo.Create(ctx, admin.ID, token, ipAddress, userAgent)
	if err != nil {
		return nil, nil, err
	}

	// Clear password before returning
	admin.Password = ""

	return admin, session, nil
}

// Logout deletes an admin session
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}

// VerifySession verifies a session token and returns the admin
func (s *AuthService) VerifySession(ctx context.Context, token string, sessionLifetime time.Duration) (*admin.Admin, error) {
	// Find session by token
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if session.IsExpired(sessionLifetime) {
		// Delete expired session
		_ = s.sessionRepo.Delete(ctx, token)
		return nil, err
	}

	// Get admin
	admin, err := s.adminRepo.FindByID(ctx, session.AdminID.String())
	if err != nil {
		return nil, err
	}

	// Update last activity if needed
	if session.ShouldRefresh() {
		_ = s.sessionRepo.UpdateLastActivity(ctx, token)
	}

	// Clear password before returning
	admin.Password = ""

	return admin, nil
}

// RefreshSession updates the last activity timestamp
func (s *AuthService) RefreshSession(ctx context.Context, token string) error {
	return s.sessionRepo.UpdateLastActivity(ctx, token)
}

// generateToken generates a secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
