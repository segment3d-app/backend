// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: assets.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createAsset = `-- name: CreateAsset :one
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
RETURNING id, uid, title, slug, "assetUrl", "assetType", "thumbnailUrl", "gaussianUrl", "pointCloudUrl", "isPrivate", status, likes, "createdAt", "updatedAt"
`

type CreateAssetParams struct {
	Uid          uuid.NullUUID `json:"uid"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Status       string        `json:"status"`
	AssetUrl     string        `json:"assetUrl"`
	AssetType    string        `json:"assetType"`
	ThumbnailUrl string        `json:"thumbnailUrl"`
	IsPrivate    bool          `json:"isPrivate"`
	Likes        int32         `json:"likes"`
}

func (q *Queries) CreateAsset(ctx context.Context, arg CreateAssetParams) (Assets, error) {
	row := q.db.QueryRowContext(ctx, createAsset,
		arg.Uid,
		arg.Title,
		arg.Slug,
		arg.Status,
		arg.AssetUrl,
		arg.AssetType,
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
		&i.AssetUrl,
		&i.AssetType,
		&i.ThumbnailUrl,
		&i.GaussianUrl,
		&i.PointCloudUrl,
		&i.IsPrivate,
		&i.Status,
		&i.Likes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
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