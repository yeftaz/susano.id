-- Drop trigger
DROP TRIGGER IF EXISTS update_admins_updated_at ON admins;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_admins_deleted_at;
DROP INDEX IF EXISTS idx_admins_role;
DROP INDEX IF EXISTS idx_admins_email;

-- Drop table
DROP TABLE IF EXISTS admins;

-- Drop enum
DROP TYPE IF EXISTS admin_role;

-- Drop wrapper function
DROP FUNCTION IF EXISTS gen_uuid_v7();

-- Note: Extensions are not dropped as they might be used by other tables
