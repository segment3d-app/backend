-- name: CreateUser :one
INSERT INTO "user" (
        username,
        email,
        password,
        phone_number,
        full_name,
        avatar
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetUserById :one
SELECT *
FROM "user"
WHERE id = $1
LIMIT 1;
-- name: GetAccountByUsername :one
SELECT *
FROM "user"
WHERE username = $1
LIMIT 1;
-- name: UpdateUser :one
UPDATE "user"
SET email = $2,
    phone_number = $3,
    full_name = $4,
    avatar = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;
-- name: UpdateUserPassword :one
UPDATE "user"
SET password = $2,
    password_change_at = now()
WHERE id = $1
RETURNING *;