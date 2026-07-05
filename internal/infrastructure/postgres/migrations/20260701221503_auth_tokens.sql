-- +goose Up

CREATE TYPE token_type AS ENUM (
    'account_activation',
    'password_reset',
    'password_reset_grant',
    'magic_link',
    'email_change'
);

CREATE TABLE tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type token_type NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_tokens_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_tokens_user ON tokens(user_id);
CREATE INDEX idx_tokens_type ON tokens(type);
CREATE INDEX idx_tokens_expires ON tokens(expires_at);

-- +goose Down

DROP TABLE tokens;

DROP TYPE token_type;