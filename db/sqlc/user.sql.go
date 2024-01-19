// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
        username,
        email,
        password,
        phone_number,
        full_name,
        avatar
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, password, phone_number, full_name, avatar, created_at, updated_at, password_changed_at
`

type CreateUserParams struct {
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	PhoneNumber sql.NullString `json:"phone_number"`
	FullName    sql.NullString `json:"full_name"`
	Avatar      sql.NullString `json:"avatar"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.PhoneNumber,
		arg.FullName,
		arg.Avatar,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.FullName,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getAccountByUsername = `-- name: GetAccountByUsername :one
SELECT id, username, email, password, phone_number, full_name, avatar, created_at, updated_at, password_changed_at
FROM "user"
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetAccountByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getAccountByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.FullName,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, email, password, phone_number, full_name, avatar, created_at, updated_at, password_changed_at
FROM "user"
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.FullName,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE "user"
SET email = $2,
    phone_number = $3,
    full_name = $4,
    avatar = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, username, email, password, phone_number, full_name, avatar, created_at, updated_at, password_changed_at
`

type UpdateUserParams struct {
	ID          uuid.UUID      `json:"id"`
	Email       string         `json:"email"`
	PhoneNumber sql.NullString `json:"phone_number"`
	FullName    sql.NullString `json:"full_name"`
	Avatar      sql.NullString `json:"avatar"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.PhoneNumber,
		arg.FullName,
		arg.Avatar,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.FullName,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE "user"
SET password = $2,
    password_changed_at = now()
WHERE id = $1
RETURNING id, username, email, password, phone_number, full_name, avatar, created_at, updated_at, password_changed_at
`

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.ID, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.FullName,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}
