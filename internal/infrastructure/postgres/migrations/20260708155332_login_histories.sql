-- +goose Up

CREATE TYPE login_status AS ENUM (
    'success',
    'failed',
    'logout',
    'expired',
    'revoked'
);

CREATE TABLE login_histories (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL,
    status login_status NOT NULL,
    ip_address INET,
    user_agent TEXT,
    device_name TEXT,
    failure_reason TEXT,
    login_at TIMESTAMP NOT NULL,
    logout_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_login_history_user ON login_histories(user_id);
CREATE INDEX idx_login_history_session ON login_histories(session_id);
CREATE INDEX idx_login_history_login_at ON login_histories(login_at DESC);
CREATE INDEX idx_login_history_status ON login_histories(status);

-- +goose Down

DROP TABLE IF EXISTS login_histories;
DROP TYPE IF EXISTS login_status;