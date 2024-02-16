-- name: CreateAsset :one
INSERT INTO "assets" (
        uid,
        title,
        slug,
        status,
        "assetUrl",
        "assetType",
        "thumbnailUrl",
        "isPrivate",
        likes
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
-- name: GetSlug :many
SELECT slug
FROM "assets"
WHERE slug LIKE $1
ORDER BY "createdAt" ASC;