// internal/auth/service/auth_service.go
package service

import (
    "context"
    "hr-backend/internal/auth/dto"
    "hr-backend/internal/auth/repository"
    "hr-backend/internal/db"
    "hr-backend/pkg/errors"
    "hr-backend/pkg/utils"
)

type IAuthService interface {
    Register(ctx context.Context, req dto.RegisterRequest) (db.User, error)
    Login(ctx context.Context, req dto.LoginRequest) (string, db.User, error)
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

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (db.User, error) {
    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return db.User{}, errors.ErrInternalServer
    }

    // Create user
    user, err := s.repo.CreateUser(ctx, req.Username, req.Email, hashedPassword)
    if err != nil {
        return db.User{}, err
    }

    return user, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, db.User, error) {
    user, err := s.repo.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return "", db.User{}, errors.ErrInvalidCredentials
    }

    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return "", db.User{}, errors.ErrInvalidCredentials
    }

    token, err := utils.GenerateToken(int(user.ID), s.jwtSecret)
    if err != nil {
        return "", db.User{}, errors.ErrInternalServer
    }

    return token, user, nil
}
