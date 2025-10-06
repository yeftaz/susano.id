package domain

import "errors"

// Common domain errors
var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSessionNotFound    = errors.New("session not found")
	ErrUnauthorized       = errors.New("unauthorized access")

	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUserDeleted        = errors.New("user account has been deleted")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password does not meet security requirements")

	// Validation errors
	ErrInvalidInput  = errors.New("invalid input data")
	ErrRequiredField = errors.New("required field is missing")

	// General errors
	ErrNotFound       = errors.New("resource not found")
	ErrForbidden      = errors.New("forbidden access")
	ErrInternalServer = errors.New("internal server error")
)
