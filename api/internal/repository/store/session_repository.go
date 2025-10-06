package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yeftaz/susano.id/api/internal/domain/store"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// Create creates a new customer session
func (r *SessionRepository) Create(ctx context.Context, customerID uuid.UUID, token, ipAddress, userAgent string) (*store.Session, error) {
	query := `
        INSERT INTO customer_sessions (id, customer_id, token, ip_address, user_agent, last_activity_at, created_at)
        VALUES (gen_uuid_v7(), $1, $2, $3, $4, NOW(), NOW())
        RETURNING id, customer_id, token, ip_address, user_agent, last_activity_at, created_at
    `

	var s store.Session
	err := r.db.QueryRowContext(ctx, query, customerID, token, ipAddress, userAgent).Scan(
		&s.ID, &s.CustomerID, &s.Token, &s.IPAddress, &s.UserAgent, &s.LastActivityAt, &s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// FindByToken retrieves a session by token
func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*store.Session, error) {
	query := `
        SELECT id, customer_id, token, ip_address, user_agent, last_activity_at, created_at
        FROM customer_sessions
        WHERE token = $1
    `

	var s store.Session
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&s.ID, &s.CustomerID, &s.Token, &s.IPAddress, &s.UserAgent, &s.LastActivityAt, &s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// UpdateLastActivity updates the last activity timestamp
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, token string) error {
	query := `
        UPDATE customer_sessions
        SET last_activity_at = NOW()
        WHERE token = $1
    `

	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// Delete deletes a session by token
func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM customer_sessions WHERE token = $1`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteByCustomerID deletes all sessions for a customer
func (r *SessionRepository) DeleteByCustomerID(ctx context.Context, customerID uuid.UUID) error {
	query := `DELETE FROM customer_sessions WHERE customer_id = $1`
	_, err := r.db.ExecContext(ctx, query, customerID)
	return err
}

// DeleteExpired deletes expired sessions
func (r *SessionRepository) DeleteExpired(ctx context.Context, sessionLifetime time.Duration) error {
	query := `
        DELETE FROM customer_sessions
        WHERE created_at < NOW() - $1::interval
    `

	_, err := r.db.ExecContext(ctx, query, sessionLifetime)
	return err
}
