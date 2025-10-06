-- Drop indexes
DROP INDEX IF EXISTS idx_password_reset_tokens_created_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;

-- Drop table
DROP TABLE IF EXISTS password_reset_tokens;
