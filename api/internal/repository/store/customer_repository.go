package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yeftaz/susano.id/api/internal/domain/store"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

// FindByEmail retrieves a customer by email
func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*store.Customer, error) {
	query := `
        SELECT id, email, password, name, avatar_path, is_active,
               email_verified_at, two_factor_secret, two_factor_recovery_codes,
               two_factor_confirmed_at, created_at, updated_at, deleted_at
        FROM customers
        WHERE email = $1 AND deleted_at IS NULL
    `

	var c store.Customer
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&c.ID, &c.Email, &c.Password, &c.Name, &c.AvatarPath,
		&c.IsActive, &c.EmailVerifiedAt, &c.TwoFactorSecret, &c.TwoFactorRecoveryCodes,
		&c.TwoFactorConfirmedAt, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// FindByID retrieves a customer by ID
func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*store.Customer, error) {
	query := `
        SELECT id, email, password, name, avatar_path, is_active,
               email_verified_at, two_factor_secret, two_factor_recovery_codes,
               two_factor_confirmed_at, created_at, updated_at, deleted_at
        FROM customers
        WHERE id = $1 AND deleted_at IS NULL
    `

	var c store.Customer
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Email, &c.Password, &c.Name, &c.AvatarPath,
		&c.IsActive, &c.EmailVerifiedAt, &c.TwoFactorSecret, &c.TwoFactorRecoveryCodes,
		&c.TwoFactorConfirmedAt, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Create creates a new customer
func (r *CustomerRepository) Create(ctx context.Context, email, passwordHash, name string) (*store.Customer, error) {
	query := `
        INSERT INTO customers (id, email, password, name, is_active, created_at, updated_at)
        VALUES (gen_uuid_v7(), $1, $2, $3, true, NOW(), NOW())
        RETURNING id, email, password, name, avatar_path, is_active,
                  email_verified_at, two_factor_secret, two_factor_recovery_codes,
                  two_factor_confirmed_at, created_at, updated_at, deleted_at
    `

	var c store.Customer
	err := r.db.QueryRowContext(ctx, query, email, passwordHash, name).Scan(
		&c.ID, &c.Email, &c.Password, &c.Name, &c.AvatarPath,
		&c.IsActive, &c.EmailVerifiedAt, &c.TwoFactorSecret, &c.TwoFactorRecoveryCodes,
		&c.TwoFactorConfirmedAt, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Update updates a customer
func (r *CustomerRepository) Update(ctx context.Context, id string, email, name *string) (*store.Customer, error) {
	// Build dynamic update query
	query := "UPDATE customers SET updated_at = NOW()"
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

	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", argCount)
	args = append(args, id)

	query += `
        RETURNING id, email, password, name, avatar_path, is_active,
                  email_verified_at, two_factor_secret, two_factor_recovery_codes,
                  two_factor_confirmed_at, created_at, updated_at, deleted_at
    `

	var c store.Customer
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.ID, &c.Email, &c.Password, &c.Name, &c.AvatarPath,
		&c.IsActive, &c.EmailVerifiedAt, &c.TwoFactorSecret, &c.TwoFactorRecoveryCodes,
		&c.TwoFactorConfirmedAt, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Delete soft deletes a customer
func (r *CustomerRepository) Delete(ctx context.Context, id string) error {
	query := `
        UPDATE customers
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

// UpdateAvatarPath updates customer avatar path
func (r *CustomerRepository) UpdateAvatarPath(ctx context.Context, id, avatarPath string) error {
	query := `
        UPDATE customers
        SET avatar_path = $1, updated_at = NOW()
        WHERE id = $2 AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, avatarPath, id)
	return err
}
