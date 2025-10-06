-- Drop indexes
DROP INDEX IF EXISTS idx_admin_sessions_last_activity_at;
DROP INDEX IF EXISTS idx_admin_sessions_created_at;
DROP INDEX IF EXISTS idx_admin_sessions_admin_id;
DROP INDEX IF EXISTS idx_admin_sessions_token;

-- Drop table
DROP TABLE IF EXISTS admin_sessions;
