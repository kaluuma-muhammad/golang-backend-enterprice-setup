-- name: CreateToken :exec

INSERT INTO tokens (
    id,
    user_id,
    type,
    token,
    expires_at,
    used_at,
    created_at,
    updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);

-- name: FindByTokenByID :one

SELECT * FROM tokens WHERE id = $1;

-- name: FindByToken :one

SELECT * FROM tokens WHERE token = $1;

-- name: MarkTokenAsUsed :exec

UPDATE tokens SET used_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: DeleteToken :exec

DELETE FROM tokens WHERE id = $1;

-- name: DeleteExpiredTokens :exec

DELETE FROM tokens WHERE expires_at < NOW();
