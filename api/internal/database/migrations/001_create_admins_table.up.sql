-- Enable UUID extension for v4 (fallback)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable UUIDv7 extension for PostgreSQL 17
CREATE EXTENSION IF NOT EXISTS "pg_uuidv7";

-- Create a wrapper function that we can easily update later
-- PostgreSQL 17: calls uuid_generate_v7() from pg_uuidv7 extension
-- PostgreSQL 18: will be updated to call built-in uuidv7()
CREATE OR REPLACE FUNCTION gen_uuid_v7()
RETURNS UUID
LANGUAGE SQL
VOLATILE
AS $$
    SELECT uuid_generate_v7();  -- PG17 with pg_uuidv7 extension
$$;

-- Create enum for admin roles
CREATE TYPE admin_role AS ENUM ('super_admin', 'admin', 'cashier');

-- Create admins table
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_path VARCHAR(500),
    role admin_role NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified_at TIMESTAMP,
    two_factor_secret TEXT,
    two_factor_recovery_codes TEXT,
    two_factor_confirmed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_admins_email ON admins(email);
CREATE INDEX idx_admins_role ON admins(role);
CREATE INDEX idx_admins_deleted_at ON admins(deleted_at) WHERE deleted_at IS NULL;

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to admins table
CREATE TRIGGER update_admins_updated_at
    BEFORE UPDATE ON admins
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
