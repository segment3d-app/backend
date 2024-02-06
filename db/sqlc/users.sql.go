// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" (
        email,
        password,
        name,
        avatar
    )
VALUES ($1, $2, $3, $4)
RETURNING uid, name, email, avatar, password, "createdAt", "updatedAt", "passwordChangedAt"
`

type CreateUserParams struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Name     sql.NullString `json:"name"`
	Avatar   sql.NullString `json:"avatar"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Password,
		arg.Name,
		arg.Avatar,
	)
	var i Users
	err := row.Scan(
		&i.Uid,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT uid, name, email, avatar, password, "createdAt", "updatedAt", "passwordChangedAt"
FROM "users"
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i Users
	err := row.Scan(
		&i.Uid,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT uid, name, email, avatar, password, "createdAt", "updatedAt", "passwordChangedAt"
FROM "users"
WHERE uid = $1
LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, uid uuid.UUID) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserById, uid)
	var i Users
	err := row.Scan(
		&i.Uid,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE "users"
SET email = $2,
    name = $3,
    avatar = $4,
    "updatedAt" = now()
WHERE uid = $1
RETURNING uid, name, email, avatar, password, "createdAt", "updatedAt", "passwordChangedAt"
`

type UpdateUserParams struct {
	Uid    uuid.UUID      `json:"uid"`
	Email  string         `json:"email"`
	Name   sql.NullString `json:"name"`
	Avatar sql.NullString `json:"avatar"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Uid,
		arg.Email,
		arg.Name,
		arg.Avatar,
	)
	var i Users
	err := row.Scan(
		&i.Uid,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE "users"
SET password = $2,
    "passwordChangedAt" = now()
WHERE uid = $1
RETURNING uid, name, email, avatar, password, "createdAt", "updatedAt", "passwordChangedAt"
`

type UpdateUserPasswordParams struct {
	Uid      uuid.UUID `json:"uid"`
	Password string    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.Uid, arg.Password)
	var i Users
	err := row.Scan(
		&i.Uid,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}
