-- name: CreateTag :one
INSERT INTO "tags" (name, slug)
VALUES ($1, $2)
RETURNING *;
-- name: GetTagsByTagsName :many
SELECT *
FROM tags
WHERE name = ANY($1);