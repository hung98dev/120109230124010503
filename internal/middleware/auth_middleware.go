// internal/middleware/auth_middleware.go
package middleware

import (
    "strings"
    "hr-backend/internal/auth/dto"
    "hr-backend/pkg/errors"
    "hr-backend/pkg/utils"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            apiErr := errors.ApiError{
                Code:    errors.ErrorCodeAuthentication,
                Message: "Authorization header is required",
            }
            c.JSON(apiErr.Code, dto.Response{
                Success: false,
                Error:   apiErr,
            })
            c.Abort()
            return
        }

        // Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            apiErr := errors.ApiError{
                Code:    errors.ErrorCodeAuthentication,
                Message: "Invalid authorization header format",
            }
            c.JSON(apiErr.Code, dto.Response{
                Success: false,
                Error:   apiErr,
            })
            c.Abort()
            return
        }

        // Validate token
        claims, err := utils.ValidateToken(parts[1], jwtSecret)
        if err != nil {
            apiErr := errors.ApiError{
                Code:    errors.ErrorCodeAuthentication,
                Message: "Invalid or expired token",
                Detail:  err.Error(),
            }
            c.JSON(apiErr.Code, dto.Response{
                Success: false,
                Error:   apiErr,
            })
            c.Abort()
            return
        }

        // Set user ID in context
        c.Set("userId", claims.UserID)
        c.Next()
    }
}