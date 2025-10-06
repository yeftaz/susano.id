-- Create customer_sessions table
CREATE TABLE customer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_customer_sessions_token ON customer_sessions(token);
CREATE INDEX idx_customer_sessions_customer_id ON customer_sessions(customer_id);
CREATE INDEX idx_customer_sessions_created_at ON customer_sessions(created_at);
CREATE INDEX idx_customer_sessions_last_activity_at ON customer_sessions(last_activity_at);
