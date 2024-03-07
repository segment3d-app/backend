-- name: CreateAssetsToTags :one
INSERT INTO "assetsToTags" ("tagsId", "assetsId")
VALUES ($1, $2)
RETURNING *;