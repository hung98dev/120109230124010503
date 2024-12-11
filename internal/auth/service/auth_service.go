// internal/auth/service/auth_service.go
package service

import (
	"context"
	"hr-backend/db"
	"hr-backend/internal/auth/dto"
	"hr-backend/internal/auth/repository"
	"hr-backend/pkg/errors"
	"hr-backend/pkg/utils"
)

type IAuthService interface {
    Register(ctx context.Context, req dto.RegisterRequest) (string, error)
    Login(ctx context.Context, req dto.LoginRequest) (string, db.GetUserByUserNameRow, error)
}

type AuthService struct {
    repo      repository.IAuthRepository
    jwtSecret string
}

func NewAuthService(repo repository.IAuthRepository, jwtSecret string) IAuthService {
    return &AuthService{
        repo:      repo,
        jwtSecret: jwtSecret,
    }
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (string, error) {
    // Hash password
    hashedPassword, err := utils.HashPassword(req.PasswordHash)
    if err != nil {
        return err.Error(), errors.ErrInternalServer
    }

    // Create user
    user, err := s.repo.CreateUser(ctx, req.Site, req.EmployeeID, hashedPassword, req.Role, req.Status, req.DepartmentID)
    if err != nil {
        return err.Error(), err
    }

    return user, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, db.GetUserByUserNameRow, error) {
    user, err := s.repo.GetUserByUserName(ctx, req.User)
    if err != nil {
        return "", db.GetUserByUserNameRow{}, errors.ErrInvalidCredentials
    }

    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return "", db.GetUserByUserNameRow{}, errors.ErrInvalidCredentials
    }

    token, err := utils.GenerateToken(int(user.ID), s.jwtSecret)
    if err != nil {
        return "", db.GetUserByUserNameRow{}, errors.ErrInternalServer
    }

    return token, user, nil
}
