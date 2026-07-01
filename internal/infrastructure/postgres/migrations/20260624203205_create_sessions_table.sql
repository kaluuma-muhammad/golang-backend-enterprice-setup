-- +goose Up

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    ip_address INET,
    device_name TEXT,
    last_used_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    revoked_reason TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    CONSTRAINT fk_sessions_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_refresh ON sessions(refresh_token);
CREATE INDEX idx_sessions_revoked ON sessions(revoked_at);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- +goose Down

DROP TABLE sessions;