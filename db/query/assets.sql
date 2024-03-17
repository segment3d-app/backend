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
-- name: GetAllAssetsByKeyword :many
SELECT a.*,
    u.name,
    u.avatar,
    u.email,
    (
        SELECT ARRAY_AGG(t.name)
        FROM "tags" AS t
            INNER JOIN "assetsToTags" AS att ON att."tagsId" = t.id
        WHERE att."assetsId" = a.id
    ) AS tag_names
FROM "assets" AS a
    LEFT JOIN "users" AS u ON u.uid = a.uid
WHERE a.title LIKE '%' || $1 || '%'
ORDER BY a."createdAt" DESC;
-- name: GetAllAssetsWithLikesInformation :many
SELECT a.*,
    u.name,
    u.avatar,
    u.email,
    CASE
        WHEN l.uid = $1 THEN TRUE
        ELSE FALSE
    END AS "isLikedByMe",
    (
        SELECT ARRAY_AGG(t.name)
        FROM "tags" AS t
            INNER JOIN "assetsToTags" AS att ON att."tagsId" = t.id
        WHERE att."assetsId" = a.id
    ) AS tag_names
FROM "assets" AS a
    LEFT JOIN "users" AS u ON u.uid = a.uid
    LEFT JOIN "likes" AS l ON l."assetsId" = a.id
    AND l.uid = $1
WHERE a.title LIKE '%' || $2 || '%'
ORDER BY a."createdAt" DESC;
-- name: GetMyAssets :many
SELECT a.*,
    CASE
        WHEN l.uid = $1 THEN TRUE
        ELSE FALSE
    END AS "isLikedByMe",
    (
        SELECT ARRAY_AGG(t.name)
        FROM "tags" AS t
            INNER JOIN "assetsToTags" AS att ON att."tagsId" = t.id
        WHERE att."assetsId" = a.id
    ) AS tag_names
FROM "assets" AS a
    LEFT JOIN "likes" AS l ON l."assetsId" = a.id
    AND l.uid = $1
WHERE a.uid = $1
    AND a.title LIKE '%' || $2 || '%'
ORDER BY "createdAt" DESC;
-- name: RemoveAsset :one
DELETE FROM "assets"
WHERE uid = $1
    AND id = $2
RETURNING *;
-- name: UpdatePointCloudUrl :one
UPDATE "assets"
SET "pointCloudUrl" = $3
WHERE uid = $1
    and id = $2
RETURNING *;
-- name: UpdateAssetUrl :one
UPDATE "assets"
SET "assetUrl" = $3
WHERE uid = $1
    and id = $2
RETURNING *;
-- name: UpdateGaussianUrl :one
UPDATE "assets"
SET "gaussianUrl" = $3,
    "status" = CASE
        WHEN "status" = 'generating splat' THEN 'completed'
        ELSE "status"
    END
WHERE uid = $1
    and id = $2
RETURNING *;
-- name: UpdateAssetStatus :one
UPDATE "assets"
SET "status" = $3
WHERE uid = $1
    and id = $2
RETURNING *;
-- name: CheckIsLiked :one
SELECT EXISTS (
        SELECT 1
        FROM "likes"
        WHERE uid = $1
            AND "assetsId" = $2
    ) AS "exists";
-- name: CreateLike :exec
INSERT INTO "likes" (uid, "assetsId")
VALUES ($1, $2);
-- name: RemoveLike :one
DELETE FROM "likes"
WHERE uid = $1
    AND "assetsId" = $2
RETURNING *;
-- name: IncreaseAssetLikes :one
UPDATE "assets"
SET likes = likes + 1
WHERE "id" = $1
RETURNING *;
-- name: DecreaseAssetLikes :one
UPDATE "assets"
SET likes = likes - 1
WHERE "id" = $1
RETURNING *;