// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (Users, error)
	GetUserByEmail(ctx context.Context, email string) (Users, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (Users, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (Users, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (Users, error)
}

var _ Querier = (*Queries)(nil)
