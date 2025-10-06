package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yeftaz/susano.id/api/internal/domain/admin"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

// FindByEmail retrieves an admin by email
func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	query := `
        SELECT id, email, password, name, avatar_path, role, is_active,
               email_verified_at, two_factor_secret, two_factor_recovery_codes,
               two_factor_confirmed_at, created_at, updated_at, deleted_at
        FROM admins
        WHERE email = $1 AND deleted_at IS NULL
    `

	var a admin.Admin
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.AvatarPath, &a.Role,
		&a.IsActive, &a.EmailVerifiedAt, &a.TwoFactorSecret, &a.TwoFactorRecoveryCodes,
		&a.TwoFactorConfirmedAt, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// FindByID retrieves an admin by ID
func (r *AdminRepository) FindByID(ctx context.Context, id string) (*admin.Admin, error) {
	query := `
        SELECT id, email, password, name, avatar_path, role, is_active,
               email_verified_at, two_factor_secret, two_factor_recovery_codes,
               two_factor_confirmed_at, created_at, updated_at, deleted_at
        FROM admins
        WHERE id = $1 AND deleted_at IS NULL
    `

	var a admin.Admin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.AvatarPath, &a.Role,
		&a.IsActive, &a.EmailVerifiedAt, &a.TwoFactorSecret, &a.TwoFactorRecoveryCodes,
		&a.TwoFactorConfirmedAt, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetAll retrieves all admins with pagination and filtering
func (r *AdminRepository) GetAll(ctx context.Context, page, limit int, search, role string) ([]*admin.Admin, int, error) {
	offset := (page - 1) * limit

	// Build query with filters
	query := `
        SELECT id, email, password, name, avatar_path, role, is_active,
               email_verified_at, two_factor_secret, two_factor_recovery_codes,
               two_factor_confirmed_at, created_at, updated_at, deleted_at
        FROM admins
        WHERE deleted_at IS NULL
    `

	countQuery := `SELECT COUNT(*) FROM admins WHERE deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argCount, argCount)
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Add role filter
	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		countQuery += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, role)
		argCount++
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get admins
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	admins := []*admin.Admin{}
	for rows.Next() {
		var a admin.Admin
		err := rows.Scan(
			&a.ID, &a.Email, &a.Password, &a.Name, &a.AvatarPath, &a.Role,
			&a.IsActive, &a.EmailVerifiedAt, &a.TwoFactorSecret, &a.TwoFactorRecoveryCodes,
			&a.TwoFactorConfirmedAt, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		admins = append(admins, &a)
	}

	return admins, total, nil
}

// Create creates a new admin
func (r *AdminRepository) Create(ctx context.Context, email, passwordHash, name, role string) (*admin.Admin, error) {
	query := `
        INSERT INTO admins (id, email, password, name, role, is_active, created_at, updated_at)
        VALUES (gen_uuid_v7(), $1, $2, $3, $4, true, NOW(), NOW())
        RETURNING id, email, password, name, avatar_path, role, is_active,
                  email_verified_at, two_factor_secret, two_factor_recovery_codes,
                  two_factor_confirmed_at, created_at, updated_at, deleted_at
    `

	var a admin.Admin
	err := r.db.QueryRowContext(ctx, query, email, passwordHash, name, role).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.AvatarPath, &a.Role,
		&a.IsActive, &a.EmailVerifiedAt, &a.TwoFactorSecret, &a.TwoFactorRecoveryCodes,
		&a.TwoFactorConfirmedAt, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// Update updates an admin
func (r *AdminRepository) Update(ctx context.Context, id string, email, name, role *string) (*admin.Admin, error) {
	// Build dynamic update query
	query := "UPDATE admins SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if email != nil {
		query += fmt.Sprintf(", email = $%d", argCount)
		args = append(args, *email)
		argCount++
	}

	if name != nil {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, *name)
		argCount++
	}

	if role != nil {
		query += fmt.Sprintf(", role = $%d", argCount)
		args = append(args, *role)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", argCount)
	args = append(args, id)

	query += `
        RETURNING id, email, password, name, avatar_path, role, is_active,
                  email_verified_at, two_factor_secret, two_factor_recovery_codes,
                  two_factor_confirmed_at, created_at, updated_at, deleted_at
    `

	var a admin.Admin
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.AvatarPath, &a.Role,
		&a.IsActive, &a.EmailVerifiedAt, &a.TwoFactorSecret, &a.TwoFactorRecoveryCodes,
		&a.TwoFactorConfirmedAt, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// Delete soft deletes an admin
func (r *AdminRepository) Delete(ctx context.Context, id string) error {
	query := `
        UPDATE admins
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `

	result, err := r.db.ExecContext(ctx, query, id)
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

// UpdateAvatarPath updates admin avatar path
func (r *AdminRepository) UpdateAvatarPath(ctx context.Context, id, avatarPath string) error {
	query := `
        UPDATE admins
        SET avatar_path = $1, updated_at = NOW()
        WHERE id = $2 AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, avatarPath, id)
	return err
}
