-- name: CreateUser :one
INSERT INTO "users" (
        email,
        password,
        name,
        avatar,
        provider
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: GetUserById :one
SELECT *
FROM "users"
WHERE uid = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM "users"
WHERE email = $1
LIMIT 1;
-- name: UpdateUser :one
UPDATE "users"
SET email = $2,
    name = $3,
    avatar = $4,
    "updatedAt" = now()
WHERE uid = $1
RETURNING *;
-- name: UpdateUserPassword :one
UPDATE "users"
SET password = $2,
    "passwordChangedAt" = now()
WHERE uid = $1
RETURNING *;