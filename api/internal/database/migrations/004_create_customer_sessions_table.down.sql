-- Drop indexes
DROP INDEX IF EXISTS idx_customer_sessions_last_activity_at;
DROP INDEX IF EXISTS idx_customer_sessions_created_at;
DROP INDEX IF EXISTS idx_customer_sessions_customer_id;
DROP INDEX IF EXISTS idx_customer_sessions_token;

-- Drop table
DROP TABLE IF EXISTS customer_sessions;
