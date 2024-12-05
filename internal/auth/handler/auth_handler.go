// internal/auth/handler/auth_handler.go
package handler

import (
    "net/http"
    "time"
    "hr-backend/internal/auth/dto"
    "hr-backend/internal/auth/service"
    "hr-backend/pkg/errors"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    service service.IAuthService
}

func NewAuthHandler(service service.IAuthService) *AuthHandler {
    return &AuthHandler{
        service: service,
    }
}

// convertTimestamptz converts pgtype.Timestamptz to time.Time
func convertTimestamptz(t interface{}) time.Time {
    switch v := t.(type) {
    case time.Time:
        return v
    default:
        return time.Time{}
    }
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        apiErr := errors.ApiError{
            Code:    errors.ErrorCodeValidation,
            Message: "Invalid request format",
            Detail:  err.Error(),
        }
        c.JSON(apiErr.Code, dto.Response{
            Success: false,
            Error:   apiErr,
        })
        return
    }

    user, err := h.service.Register(c.Request.Context(), req)
    if err != nil {
        apiErr := errors.NewApiError(err)
        c.JSON(apiErr.Code, dto.Response{
            Success: false,
            Error:   apiErr,
        })
        return
    }

    // Convert timestamps
    createdAt := convertTimestamptz(user.CreatedAt)
    updatedAt := convertTimestamptz(user.UpdatedAt)

    c.JSON(http.StatusCreated, dto.Response{
        Success: true,
        Data: dto.UserResponse{
            ID:        int64(user.ID), // Convert int32 to int64
            Username:  user.Username,
            Email:     user.Email,
            CreatedAt: createdAt.Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt: updatedAt.Format("2006-01-02T15:04:05Z07:00"),
        },
    })
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        apiErr := errors.ApiError{
            Code:    errors.ErrorCodeValidation,
            Message: "Invalid request format",
            Detail:  err.Error(),
        }
        c.JSON(apiErr.Code, dto.Response{
            Success: false,
            Error:   apiErr,
        })
        return
    }

    token, user, err := h.service.Login(c.Request.Context(), req)
    if err != nil {
        apiErr := errors.NewApiError(err)
        c.JSON(apiErr.Code, dto.Response{
            Success: false,
            Error:   apiErr,
        })
        return
    }

    // Convert timestamps
    createdAt := convertTimestamptz(user.CreatedAt)
    updatedAt := convertTimestamptz(user.UpdatedAt)

    c.JSON(http.StatusOK, dto.Response{
        Success: true,
        Data: dto.LoginResponse{
            Token: token,
            User: dto.UserResponse{
                ID:        int64(user.ID),
                Username:  user.Username,
                Email:     user.Email,
                CreatedAt: createdAt.Format("2006-01-02T15:04:05Z07:00"),
                UpdatedAt: updatedAt.Format("2006-01-02T15:04:05Z07:00"),
            },
        },
    })
}
