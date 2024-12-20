// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  full_name,
  password
) VALUES (
  $1, $2, $3, $4
) RETURNING username, email, full_name, password, password_changed_at, created_at, is_email_verified, role
`

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, email, full_name, password, password_changed_at, created_at, is_email_verified, role FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET 
  password = COALESCE($1, password),
  password_changed_at = COALESCE($2, password_changed_at),
  full_name = COALESCE($3, full_name),
  email = COALESCE($4, email),
  is_email_verified = COALESCE($5, is_email_verified)
WHERE username = $6 RETURNING username, email, full_name, password, password_changed_at, created_at, is_email_verified, role
`

type UpdateUserParams struct {
	Password          pgtype.Text        `json:"password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	FullName          pgtype.Text        `json:"full_name"`
	Email             pgtype.Text        `json:"email"`
	IsEmailVerified   pgtype.Bool        `json:"is_email_verified"`
	Username          string             `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Password,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
		&i.Role,
	)
	return i, err
}
