-- name: CreateUser :exec

INSERT INTO users (
    id,
    email,
    password,
    first_name,
    last_name,
    is_verified,
    created_at,
    updated_at
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8
);

-- name: GetUserByID :one

SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one

SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    is_verified = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;