// internal/auth/repository/auth_repository.go
package repository

import (
	"context"
	"hr-backend/db"
	"hr-backend/pkg/errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IAuthRepository định nghĩa interface cho repository
type IAuthRepository interface {
	CreateUser(ctx context.Context, site pgtype.Text, employee_id, passwordHash string, role, status, department_id pgtype.Int2) (string, error)
	GetUserByUserName(ctx context.Context, username string) (db.GetUserByUserNameRow, error)
}

// AuthRepository struct implement IAuthRepository interface
type AuthRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewAuthRepository tạo instance mới của AuthRepository
func NewAuthRepository(pool *pgxpool.Pool) IAuthRepository {
	return &AuthRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// handleDuplicateError xử lý lỗi duplicate key
func handleDuplicateError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "users_email_key") {
		return errors.ErrEmailExists
	}
	if strings.Contains(errMsg, "users_username_key") {
		return errors.ErrUsernameExists
	}
	return errors.ErrDatabaseError
}

// CreateUser implements IAuthRepository
func (r *AuthRepository) CreateUser(ctx context.Context, site pgtype.Text, employee_id, passwordHash string, role, status, department_id pgtype.Int2) (string, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Site:         site,
		EmployeeID:   employee_id,
		Username:     employee_id,
		PasswordHash: passwordHash,
		DepartmentID: department_id,
		Role:         role,
		Status:       status,
	})

	if err != nil {
		return "ERR", handleDuplicateError(err)
	}

	return user, nil
}

// GetUserBy UserName OR EmployeeID implements IAuthRepository
func (r *AuthRepository) GetUserByUserName(ctx context.Context, username string) (db.GetUserByUserNameRow, error) {
	user, err := r.queries.GetUserByUserName(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return db.GetUserByUserNameRow{}, errors.ErrUserNotFound
		}
		return db.GetUserByUserNameRow{}, errors.ErrDatabaseError
	}
	return user, nil
}
