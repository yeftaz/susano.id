package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/yeftaz/susano.id/api/internal/domain/store"
	storeRepo "github.com/yeftaz/susano.id/api/internal/repository/store"
)

type AuthService struct {
	customerRepo *storeRepo.CustomerRepository
	sessionRepo  *storeRepo.SessionRepository
}

func NewAuthService(customerRepo *storeRepo.CustomerRepository, sessionRepo *storeRepo.SessionRepository) *AuthService {
	return &AuthService{
		customerRepo: customerRepo,
		sessionRepo:  sessionRepo,
	}
}

// Login authenticates a customer and creates a session
func (s *AuthService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*store.Customer, *store.Session, error) {
	// Find customer by email
	customer, err := s.customerRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password)); err != nil {
		return nil, nil, err
	}

	// Check if customer can make purchases
	if !customer.CanPurchase() {
		return nil, nil, err
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, nil, err
	}

	// Create session
	session, err := s.sessionRepo.Create(ctx, customer.ID, token, ipAddress, userAgent)
	if err != nil {
		return nil, nil, err
	}

	// Clear password before returning
	customer.Password = ""

	return customer, session, nil
}

// Register creates a new customer account
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*store.Customer, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create customer
	customer, err := s.customerRepo.Create(ctx, email, string(hashedPassword), name)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	customer.Password = ""

	return customer, nil
}

// Logout deletes a customer session
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}

// VerifySession verifies a session token and returns the customer
func (s *AuthService) VerifySession(ctx context.Context, token string, sessionLifetime time.Duration) (*store.Customer, error) {
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

	// Get customer
	customer, err := s.customerRepo.FindByID(ctx, session.CustomerID.String())
	if err != nil {
		return nil, err
	}

	// Update last activity if needed
	if session.ShouldRefresh() {
		_ = s.sessionRepo.UpdateLastActivity(ctx, token)
	}

	// Clear password before returning
	customer.Password = ""

	return customer, nil
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
