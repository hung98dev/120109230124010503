// internal/auth/dto/auth_dto.go
package dto

import "github.com/jackc/pgx/v5/pgtype"

type RegisterRequest struct {
	Site         pgtype.Text `json:"site"`
	EmployeeID   string      `json:"employee_id" binding:"required,min=1,max=20"`
	Username     string      `json:"username" binding:"min=6,max=20"`
	PasswordHash string      `json:"password_hash" binding:"required,min=6"`
	DepartmentID pgtype.Int2 `json:"department_id" binding:"required"`
	Role         pgtype.Int2 `json:"role" binding:"required"`
	Status       pgtype.Int2 `json:"status" binding:"required"`
}

type LoginRequest struct {
	User     string `json:"user" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID           int32       `json:"id"`
	EmployeeID   string      `json:"employee_id"`
	Username     string      `json:"username"`
	Email        pgtype.Text `json:"email"`
	DepartmentID pgtype.Int2 `json:"department_id"`
	Role         pgtype.Int2 `json:"role"`
	Status       pgtype.Int2 `json:"status"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
