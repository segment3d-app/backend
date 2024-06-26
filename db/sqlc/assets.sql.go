// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: assets.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const checkIsLiked = `-- name: CheckIsLiked :one
SELECT EXISTS (
        SELECT 1
        FROM "likes"
        WHERE uid = $1
            AND "assetsId" = $2
    ) AS "exists"
`

type CheckIsLikedParams struct {
	Uid      uuid.UUID `json:"uid"`
	AssetsId uuid.UUID `json:"assetsId"`
}

func (q *Queries) CheckIsLiked(ctx context.Context, arg CheckIsLikedParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkIsLiked, arg.Uid, arg.AssetsId)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createAsset = `-- name: CreateAsset :one
INSERT INTO "assets" (
        uid,
        title,
        slug,
        status,
        "photoDirUrl",
        "type",
        "thumbnailUrl",
        "isPrivate",
        likes
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type CreateAssetParams struct {
	Uid          uuid.UUID `json:"uid"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Status       string    `json:"status"`
	PhotoDirUrl  string    `json:"photoDirUrl"`
	Type         string    `json:"type"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
	IsPrivate    bool      `json:"isPrivate"`
	Likes        int32     `json:"likes"`
}

func (q *Queries) CreateAsset(ctx context.Context, arg CreateAssetParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, createAsset,
		arg.Uid,
		arg.Title,
		arg.Slug,
		arg.Status,
		arg.PhotoDirUrl,
		arg.Type,
		arg.ThumbnailUrl,
		arg.IsPrivate,
		arg.Likes,
	)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createLike = `-- name: CreateLike :exec
INSERT INTO "likes" (uid, "assetsId")
VALUES ($1, $2)
`

type CreateLikeParams struct {
	Uid      uuid.UUID `json:"uid"`
	AssetsId uuid.UUID `json:"assetsId"`
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) error {
	_, err := q.db.ExecContext(ctx, createLike, arg.Uid, arg.AssetsId)
	return err
}

const decreaseAssetLikes = `-- name: DecreaseAssetLikes :one
UPDATE "assets"
SET likes = likes - 1
WHERE "id" = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

func (q *Queries) DecreaseAssetLikes(ctx context.Context, id uuid.UUID) (Assets, error) {
	row := q.db.QueryRowContext(ctx, decreaseAssetLikes, id)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllAssets = `-- name: GetAllAssets :many
SELECT a.id, a.uid, a.title, a.slug, a.type, a."thumbnailUrl", a."photoDirUrl", a."splatUrl", a."pclUrl", a."pclColmapUrl", a."segmentedPclDirUrl", a."segmentedSplatDirUrl", a."isPrivate", a.status, a.likes, a."createdAt", a."updatedAt",
    u.name,
    u.avatar,
    u.email
FROM "assets" AS a
    LEFT JOIN "users" AS u ON u.uid = a.uid
ORDER BY a."createdAt" DESC
`

type GetAllAssetsRow struct {
	ID                   uuid.UUID      `json:"id"`
	Uid                  uuid.UUID      `json:"uid"`
	Title                string         `json:"title"`
	Slug                 string         `json:"slug"`
	Type                 string         `json:"type"`
	ThumbnailUrl         string         `json:"thumbnailUrl"`
	PhotoDirUrl          string         `json:"photoDirUrl"`
	SplatUrl             sql.NullString `json:"splatUrl"`
	PclUrl               sql.NullString `json:"pclUrl"`
	PclColmapUrl         sql.NullString `json:"pclColmapUrl"`
	SegmentedPclDirUrl   sql.NullString `json:"segmentedPclDirUrl"`
	SegmentedSplatDirUrl sql.NullString `json:"segmentedSplatDirUrl"`
	IsPrivate            bool           `json:"isPrivate"`
	Status               string         `json:"status"`
	Likes                int32          `json:"likes"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	Name                 sql.NullString `json:"name"`
	Avatar               sql.NullString `json:"avatar"`
	Email                sql.NullString `json:"email"`
}

func (q *Queries) GetAllAssets(ctx context.Context) ([]GetAllAssetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAssets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllAssetsRow{}
	for rows.Next() {
		var i GetAllAssetsRow
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.Title,
			&i.Slug,
			&i.Type,
			&i.ThumbnailUrl,
			&i.PhotoDirUrl,
			&i.SplatUrl,
			&i.PclUrl,
			&i.PclColmapUrl,
			&i.SegmentedPclDirUrl,
			&i.SegmentedSplatDirUrl,
			&i.IsPrivate,
			&i.Status,
			&i.Likes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Avatar,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllAssetsByKeyword = `-- name: GetAllAssetsByKeyword :many
SELECT a.id, a.uid, a.title, a.slug, a.type, a."thumbnailUrl", a."photoDirUrl", a."splatUrl", a."pclUrl", a."pclColmapUrl", a."segmentedPclDirUrl", a."segmentedSplatDirUrl", a."isPrivate", a.status, a.likes, a."createdAt", a."updatedAt",
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
WHERE a.title LIKE '%' || $1 || '%' and a."isPrivate" = false
ORDER BY a."createdAt" DESC
`

type GetAllAssetsByKeywordRow struct {
	ID                   uuid.UUID      `json:"id"`
	Uid                  uuid.UUID      `json:"uid"`
	Title                string         `json:"title"`
	Slug                 string         `json:"slug"`
	Type                 string         `json:"type"`
	ThumbnailUrl         string         `json:"thumbnailUrl"`
	PhotoDirUrl          string         `json:"photoDirUrl"`
	SplatUrl             sql.NullString `json:"splatUrl"`
	PclUrl               sql.NullString `json:"pclUrl"`
	PclColmapUrl         sql.NullString `json:"pclColmapUrl"`
	SegmentedPclDirUrl   sql.NullString `json:"segmentedPclDirUrl"`
	SegmentedSplatDirUrl sql.NullString `json:"segmentedSplatDirUrl"`
	IsPrivate            bool           `json:"isPrivate"`
	Status               string         `json:"status"`
	Likes                int32          `json:"likes"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	Name                 sql.NullString `json:"name"`
	Avatar               sql.NullString `json:"avatar"`
	Email                sql.NullString `json:"email"`
	TagNames             []string       `json:"tag_names"`
}

func (q *Queries) GetAllAssetsByKeyword(ctx context.Context, dollar_1 sql.NullString) ([]GetAllAssetsByKeywordRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAssetsByKeyword, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllAssetsByKeywordRow{}
	for rows.Next() {
		var i GetAllAssetsByKeywordRow
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.Title,
			&i.Slug,
			&i.Type,
			&i.ThumbnailUrl,
			&i.PhotoDirUrl,
			&i.SplatUrl,
			&i.PclUrl,
			&i.PclColmapUrl,
			&i.SegmentedPclDirUrl,
			&i.SegmentedSplatDirUrl,
			&i.IsPrivate,
			&i.Status,
			&i.Likes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Avatar,
			&i.Email,
			pq.Array(&i.TagNames),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllAssetsWithLikesInformation = `-- name: GetAllAssetsWithLikesInformation :many
SELECT a.id, a.uid, a.title, a.slug, a.type, a."thumbnailUrl", a."photoDirUrl", a."splatUrl", a."pclUrl", a."pclColmapUrl", a."segmentedPclDirUrl", a."segmentedSplatDirUrl", a."isPrivate", a.status, a.likes, a."createdAt", a."updatedAt",
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
WHERE a.title LIKE '%' || $2 || '%' and a."isPrivate" = false
ORDER BY a."createdAt" DESC
`

type GetAllAssetsWithLikesInformationParams struct {
	Uid     uuid.UUID      `json:"uid"`
	Column2 sql.NullString `json:"column_2"`
}

type GetAllAssetsWithLikesInformationRow struct {
	ID                   uuid.UUID      `json:"id"`
	Uid                  uuid.UUID      `json:"uid"`
	Title                string         `json:"title"`
	Slug                 string         `json:"slug"`
	Type                 string         `json:"type"`
	ThumbnailUrl         string         `json:"thumbnailUrl"`
	PhotoDirUrl          string         `json:"photoDirUrl"`
	SplatUrl             sql.NullString `json:"splatUrl"`
	PclUrl               sql.NullString `json:"pclUrl"`
	PclColmapUrl         sql.NullString `json:"pclColmapUrl"`
	SegmentedPclDirUrl   sql.NullString `json:"segmentedPclDirUrl"`
	SegmentedSplatDirUrl sql.NullString `json:"segmentedSplatDirUrl"`
	IsPrivate            bool           `json:"isPrivate"`
	Status               string         `json:"status"`
	Likes                int32          `json:"likes"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	Name                 sql.NullString `json:"name"`
	Avatar               sql.NullString `json:"avatar"`
	Email                sql.NullString `json:"email"`
	IsLikedByMe          bool           `json:"isLikedByMe"`
	TagNames             []string       `json:"tag_names"`
}

func (q *Queries) GetAllAssetsWithLikesInformation(ctx context.Context, arg GetAllAssetsWithLikesInformationParams) ([]GetAllAssetsWithLikesInformationRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAssetsWithLikesInformation, arg.Uid, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllAssetsWithLikesInformationRow{}
	for rows.Next() {
		var i GetAllAssetsWithLikesInformationRow
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.Title,
			&i.Slug,
			&i.Type,
			&i.ThumbnailUrl,
			&i.PhotoDirUrl,
			&i.SplatUrl,
			&i.PclUrl,
			&i.PclColmapUrl,
			&i.SegmentedPclDirUrl,
			&i.SegmentedSplatDirUrl,
			&i.IsPrivate,
			&i.Status,
			&i.Likes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Avatar,
			&i.Email,
			&i.IsLikedByMe,
			pq.Array(&i.TagNames),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAssetsById = `-- name: GetAssetsById :one
SELECT id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
FROM "assets"
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAssetsById(ctx context.Context, id uuid.UUID) (Assets, error) {
	row := q.db.QueryRowContext(ctx, getAssetsById, id)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAssetsBySlug = `-- name: GetAssetsBySlug :one
SELECT id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
FROM "assets"
WHERE slug = $1
LIMIT 1
`

func (q *Queries) GetAssetsBySlug(ctx context.Context, slug string) (Assets, error) {
	row := q.db.QueryRowContext(ctx, getAssetsBySlug, slug)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAssetsByUid = `-- name: GetAssetsByUid :many
SELECT id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
FROM "assets"
WHERE uid = $1
ORDER BY "createdAt" DESC
`

func (q *Queries) GetAssetsByUid(ctx context.Context, uid uuid.UUID) ([]Assets, error) {
	rows, err := q.db.QueryContext(ctx, getAssetsByUid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Assets{}
	for rows.Next() {
		var i Assets
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.Title,
			&i.Slug,
			&i.Type,
			&i.ThumbnailUrl,
			&i.PhotoDirUrl,
			&i.SplatUrl,
			&i.PclUrl,
			&i.PclColmapUrl,
			&i.SegmentedPclDirUrl,
			&i.SegmentedSplatDirUrl,
			&i.IsPrivate,
			&i.Status,
			&i.Likes,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMyAssets = `-- name: GetMyAssets :many
SELECT a.id, a.uid, a.title, a.slug, a.type, a."thumbnailUrl", a."photoDirUrl", a."splatUrl", a."pclUrl", a."pclColmapUrl", a."segmentedPclDirUrl", a."segmentedSplatDirUrl", a."isPrivate", a.status, a.likes, a."createdAt", a."updatedAt",
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
ORDER BY "createdAt" DESC
`

type GetMyAssetsParams struct {
	Uid     uuid.UUID      `json:"uid"`
	Column2 sql.NullString `json:"column_2"`
}

type GetMyAssetsRow struct {
	ID                   uuid.UUID      `json:"id"`
	Uid                  uuid.UUID      `json:"uid"`
	Title                string         `json:"title"`
	Slug                 string         `json:"slug"`
	Type                 string         `json:"type"`
	ThumbnailUrl         string         `json:"thumbnailUrl"`
	PhotoDirUrl          string         `json:"photoDirUrl"`
	SplatUrl             sql.NullString `json:"splatUrl"`
	PclUrl               sql.NullString `json:"pclUrl"`
	PclColmapUrl         sql.NullString `json:"pclColmapUrl"`
	SegmentedPclDirUrl   sql.NullString `json:"segmentedPclDirUrl"`
	SegmentedSplatDirUrl sql.NullString `json:"segmentedSplatDirUrl"`
	IsPrivate            bool           `json:"isPrivate"`
	Status               string         `json:"status"`
	Likes                int32          `json:"likes"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	IsLikedByMe          sql.NullBool   `json:"isLikedByMe"`
	TagNames             []string       `json:"tag_names"`
}

func (q *Queries) GetMyAssets(ctx context.Context, arg GetMyAssetsParams) ([]GetMyAssetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMyAssets, arg.Uid, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMyAssetsRow{}
	for rows.Next() {
		var i GetMyAssetsRow
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.Title,
			&i.Slug,
			&i.Type,
			&i.ThumbnailUrl,
			&i.PhotoDirUrl,
			&i.SplatUrl,
			&i.PclUrl,
			&i.PclColmapUrl,
			&i.SegmentedPclDirUrl,
			&i.SegmentedSplatDirUrl,
			&i.IsPrivate,
			&i.Status,
			&i.Likes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsLikedByMe,
			pq.Array(&i.TagNames),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSlug = `-- name: GetSlug :many
SELECT slug
FROM "assets"
WHERE slug LIKE $1
ORDER BY "createdAt" ASC
`

func (q *Queries) GetSlug(ctx context.Context, slug string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getSlug, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		items = append(items, slug)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const increaseAssetLikes = `-- name: IncreaseAssetLikes :one
UPDATE "assets"
SET likes = likes + 1
WHERE "id" = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

func (q *Queries) IncreaseAssetLikes(ctx context.Context, id uuid.UUID) (Assets, error) {
	row := q.db.QueryRowContext(ctx, increaseAssetLikes, id)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeAsset = `-- name: RemoveAsset :one
DELETE FROM "assets"
WHERE uid = $1
    AND id = $2
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type RemoveAssetParams struct {
	Uid uuid.UUID `json:"uid"`
	ID  uuid.UUID `json:"id"`
}

func (q *Queries) RemoveAsset(ctx context.Context, arg RemoveAssetParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, removeAsset, arg.Uid, arg.ID)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeLike = `-- name: RemoveLike :one
DELETE FROM "likes"
WHERE uid = $1
    AND "assetsId" = $2
RETURNING uid, "assetsId", "createdAt", "updatedAt"
`

type RemoveLikeParams struct {
	Uid      uuid.UUID `json:"uid"`
	AssetsId uuid.UUID `json:"assetsId"`
}

func (q *Queries) RemoveLike(ctx context.Context, arg RemoveLikeParams) (Likes, error) {
	row := q.db.QueryRowContext(ctx, removeLike, arg.Uid, arg.AssetsId)
	var i Likes
	err := row.Scan(
		&i.Uid,
		&i.AssetsId,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAssetStatus = `-- name: UpdateAssetStatus :one
UPDATE "assets"
SET "status" = $3
WHERE uid = $1
    and id = $2
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdateAssetStatusParams struct {
	Uid    uuid.UUID `json:"uid"`
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (q *Queries) UpdateAssetStatus(ctx context.Context, arg UpdateAssetStatusParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updateAssetStatus, arg.Uid, arg.ID, arg.Status)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePTvUrl = `-- name: UpdatePTvUrl :one
UPDATE "assets"
SET "segmentedPclDirUrl" = $2
WHERE id = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdatePTvUrlParams struct {
	ID                 uuid.UUID      `json:"id"`
	SegmentedPclDirUrl sql.NullString `json:"segmentedPclDirUrl"`
}

func (q *Queries) UpdatePTvUrl(ctx context.Context, arg UpdatePTvUrlParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updatePTvUrl, arg.ID, arg.SegmentedPclDirUrl)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePointCloudUrlFromColmap = `-- name: UpdatePointCloudUrlFromColmap :one
UPDATE "assets"
SET "pclColmapUrl" = $2
WHERE id = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdatePointCloudUrlFromColmapParams struct {
	ID           uuid.UUID      `json:"id"`
	PclColmapUrl sql.NullString `json:"pclColmapUrl"`
}

func (q *Queries) UpdatePointCloudUrlFromColmap(ctx context.Context, arg UpdatePointCloudUrlFromColmapParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updatePointCloudUrlFromColmap, arg.ID, arg.PclColmapUrl)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePointCloudUrlFromLidar = `-- name: UpdatePointCloudUrlFromLidar :one
UPDATE "assets"
SET "pclUrl" = $3
WHERE uid = $1
    and id = $2
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdatePointCloudUrlFromLidarParams struct {
	Uid    uuid.UUID      `json:"uid"`
	ID     uuid.UUID      `json:"id"`
	PclUrl sql.NullString `json:"pclUrl"`
}

func (q *Queries) UpdatePointCloudUrlFromLidar(ctx context.Context, arg UpdatePointCloudUrlFromLidarParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updatePointCloudUrlFromLidar, arg.Uid, arg.ID, arg.PclUrl)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateSagaUrl = `-- name: UpdateSagaUrl :one
UPDATE "assets"
SET "segmentedSplatDirUrl" = $2
WHERE id = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdateSagaUrlParams struct {
	ID                   uuid.UUID      `json:"id"`
	SegmentedSplatDirUrl sql.NullString `json:"segmentedSplatDirUrl"`
}

func (q *Queries) UpdateSagaUrl(ctx context.Context, arg UpdateSagaUrlParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updateSagaUrl, arg.ID, arg.SegmentedSplatDirUrl)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateSplatUrl = `-- name: UpdateSplatUrl :one
UPDATE "assets"
SET "splatUrl" = $2
WHERE id = $1
RETURNING id, uid, title, slug, type, "thumbnailUrl", "photoDirUrl", "splatUrl", "pclUrl", "pclColmapUrl", "segmentedPclDirUrl", "segmentedSplatDirUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type UpdateSplatUrlParams struct {
	ID       uuid.UUID      `json:"id"`
	SplatUrl sql.NullString `json:"splatUrl"`
}

func (q *Queries) UpdateSplatUrl(ctx context.Context, arg UpdateSplatUrlParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, updateSplatUrl, arg.ID, arg.SplatUrl)
	var i Assets
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Slug,
		&i.Type,
		&i.ThumbnailUrl,
		&i.PhotoDirUrl,
		&i.SplatUrl,
		&i.PclUrl,
		&i.PclColmapUrl,
		&i.SegmentedPclDirUrl,
		&i.SegmentedSplatDirUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
