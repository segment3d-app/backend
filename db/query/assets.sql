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
-- name: GetAssetsByUid :many
SELECT *
FROM "assets"
WHERE uid = $1
ORDER BY "createdAt" DESC;
-- name: GetAssetsBySlug :one
SELECT *
FROM "assets"
WHERE slug = $1
LIMIT 1;
-- name: GetAssetsById :one
SELECT *
FROM "assets"
WHERE id = $1
LIMIT 1;
-- name: GetAllAssets :many
SELECT a.*,
    u.name,
    u.avatar,
    u.email
FROM "assets" AS a
    LEFT JOIN "users" AS u ON u.uid = a.uid
ORDER BY a."createdAt" DESC;
-- name: GetMyAssets :many
SELECT *
FROM "assets"
WHERE uid = $1
ORDER BY "createdAt" DESC;
-- name: RemoveAsset :one
DELETE FROM "assets"
WHERE uid = $1 AND id = $2
RETURNING *;