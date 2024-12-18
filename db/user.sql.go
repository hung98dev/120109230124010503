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
INSERT INTO public.user (
    site,
    employee_id,
    username,
    password_hash,
    department_id,
    role,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING employee_id
`

type CreateUserParams struct {
	Site         pgtype.Text `json:"site"`
	EmployeeID   string      `json:"employee_id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"password_hash"`
	DepartmentID pgtype.Int2 `json:"department_id"`
	Role         pgtype.Int2 `json:"role"`
	Status       pgtype.Int2 `json:"status"`
}

// sqlc/queries/user.sql
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (string, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Site,
		arg.EmployeeID,
		arg.Username,
		arg.PasswordHash,
		arg.DepartmentID,
		arg.Role,
		arg.Status,
	)
	var employee_id string
	err := row.Scan(&employee_id)
	return employee_id, err
}

const getUserByEmployeeID = `-- name: GetUserByEmployeeID :one
SELECT 
    id,
    site,
    employee_id,
    username,
    password_hash,
    email,
    phone,
    department_id,
    role,
    status,
    status_reason FROM public.user
WHERE employee_id = $1 LIMIT 1
`

type GetUserByEmployeeIDRow struct {
	ID           int32       `json:"id"`
	Site         pgtype.Text `json:"site"`
	EmployeeID   string      `json:"employee_id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"password_hash"`
	Email        pgtype.Text `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	DepartmentID pgtype.Int2 `json:"department_id"`
	Role         pgtype.Int2 `json:"role"`
	Status       pgtype.Int2 `json:"status"`
	StatusReason pgtype.Text `json:"status_reason"`
}

func (q *Queries) GetUserByEmployeeID(ctx context.Context, employeeID string) (GetUserByEmployeeIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmployeeID, employeeID)
	var i GetUserByEmployeeIDRow
	err := row.Scan(
		&i.ID,
		&i.Site,
		&i.EmployeeID,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.Phone,
		&i.DepartmentID,
		&i.Role,
		&i.Status,
		&i.StatusReason,
	)
	return i, err
}

const getUserByUserName = `-- name: GetUserByUserName :one
SELECT 
    id,
    site,
    employee_id,
    username,
    password_hash,
    email,
    phone,
    department_id,
    role,
    status,
    status_reason
FROM public.user
WHERE username = $1 LIMIT 1
`

type GetUserByUserNameRow struct {
	ID           int32       `json:"id"`
	Site         pgtype.Text `json:"site"`
	EmployeeID   string      `json:"employee_id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"password_hash"`
	Email        pgtype.Text `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	DepartmentID pgtype.Int2 `json:"department_id"`
	Role         pgtype.Int2 `json:"role"`
	Status       pgtype.Int2 `json:"status"`
	StatusReason pgtype.Text `json:"status_reason"`
}

func (q *Queries) GetUserByUserName(ctx context.Context, username string) (GetUserByUserNameRow, error) {
	row := q.db.QueryRow(ctx, getUserByUserName, username)
	var i GetUserByUserNameRow
	err := row.Scan(
		&i.ID,
		&i.Site,
		&i.EmployeeID,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.Phone,
		&i.DepartmentID,
		&i.Role,
		&i.Status,
		&i.StatusReason,
	)
	return i, err
}
