package admin

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yeftaz/susano.id/api/internal/domain/admin"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// Create creates a new admin session
func (r *SessionRepository) Create(ctx context.Context, adminID uuid.UUID, token, ipAddress, userAgent string) (*admin.Session, error) {
	query := `
        INSERT INTO admin_sessions (id, admin_id, token, ip_address, user_agent, last_activity_at, created_at)
        VALUES (gen_uuid_v7(), $1, $2, $3, $4, NOW(), NOW())
        RETURNING id, admin_id, token, ip_address, user_agent, last_activity_at, created_at
    `

	var s admin.Session
	err := r.db.QueryRowContext(ctx, query, adminID, token, ipAddress, userAgent).Scan(
		&s.ID, &s.AdminID, &s.Token, &s.IPAddress, &s.UserAgent, &s.LastActivityAt, &s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// FindByToken retrieves a session by token
func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*admin.Session, error) {
	query := `
        SELECT id, admin_id, token, ip_address, user_agent, last_activity_at, created_at
        FROM admin_sessions
        WHERE token = $1
    `

	var s admin.Session
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&s.ID, &s.AdminID, &s.Token, &s.IPAddress, &s.UserAgent, &s.LastActivityAt, &s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// UpdateLastActivity updates the last activity timestamp
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, token string) error {
	query := `
        UPDATE admin_sessions
        SET last_activity_at = NOW()
        WHERE token = $1
    `

	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// Delete deletes a session by token
func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM admin_sessions WHERE token = $1`

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

// DeleteByAdminID deletes all sessions for an admin
func (r *SessionRepository) DeleteByAdminID(ctx context.Context, adminID uuid.UUID) error {
	query := `DELETE FROM admin_sessions WHERE admin_id = $1`
	_, err := r.db.ExecContext(ctx, query, adminID)
	return err
}

// DeleteExpired deletes expired sessions
func (r *SessionRepository) DeleteExpired(ctx context.Context, sessionLifetime time.Duration) error {
	query := `
        DELETE FROM admin_sessions
        WHERE created_at < NOW() - $1::interval
    `

	_, err := r.db.ExecContext(ctx, query, sessionLifetime)
	return err
}
