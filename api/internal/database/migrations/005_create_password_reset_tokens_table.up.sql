-- Create password_reset_tokens table
CREATE TABLE password_reset_tokens (
    email VARCHAR(255) PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index for token lookup
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);

-- Create index for cleanup (delete old tokens)
CREATE INDEX idx_password_reset_tokens_created_at ON password_reset_tokens(created_at);
