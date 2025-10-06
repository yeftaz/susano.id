-- Create admin_sessions table
CREATE TABLE admin_sessions (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    admin_id UUID NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_admin_sessions_token ON admin_sessions(token);
CREATE INDEX idx_admin_sessions_admin_id ON admin_sessions(admin_id);
CREATE INDEX idx_admin_sessions_created_at ON admin_sessions(created_at);
CREATE INDEX idx_admin_sessions_last_activity_at ON admin_sessions(last_activity_at);
