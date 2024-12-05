// internal/auth/repository/auth_repository.go
package repository

import (
    "context"
    "strings"
    "hr-backend/internal/db"
    "hr-backend/pkg/errors"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// IAuthRepository định nghĩa interface cho repository
type IAuthRepository interface {
    CreateUser(ctx context.Context, username, email, passwordHash string) (db.User, error)
    GetUserByEmail(ctx context.Context, email string) (db.User, error)
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
func (r *AuthRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (db.User, error) {
    user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Username:     username,
        Email:        email,
        PasswordHash: passwordHash,
    })

    if err != nil {
        return db.User{}, handleDuplicateError(err)
    }

    return user, nil
}

// GetUserByEmail implements IAuthRepository
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
    user, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        if err == pgx.ErrNoRows {
            return db.User{}, errors.ErrUserNotFound
        }
        return db.User{}, errors.ErrDatabaseError
    }
    return user, nil
}