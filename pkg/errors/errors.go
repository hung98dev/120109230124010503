// pkg/errors/errors.go
package errors

import "errors"

// Custom error types
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound      = errors.New("user not found")
    ErrUserExists        = errors.New("user already exists")
    ErrEmailExists       = errors.New("email already exists")
    ErrUsernameExists    = errors.New("username already exists")
    ErrInvalidEmail      = errors.New("invalid email format")
    ErrInvalidUsername   = errors.New("invalid username format")
    ErrInvalidPassword   = errors.New("invalid password format")
    ErrInternalServer    = errors.New("internal server error")
    ErrDatabaseError     = errors.New("database error")
)

// ApiError represents a custom error response
type ApiError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// Error codes
const (
    ErrorCodeValidation     = 400
    ErrorCodeAuthentication = 401
    ErrorCodeForbidden      = 403
    ErrorCodeNotFound       = 404
    ErrorCodeConflict       = 409
    ErrorCodeInternal       = 500
)

// Error mapping
func NewApiError(err error) ApiError {
    switch err {
    case ErrInvalidCredentials:
        return ApiError{
            Code:    ErrorCodeAuthentication,
            Message: "Authentication failed",
            Detail:  err.Error(),
        }
    case ErrUserNotFound:
        return ApiError{
            Code:    ErrorCodeNotFound,
            Message: "Resource not found",
            Detail:  err.Error(),
        }
    case ErrUserExists, ErrEmailExists, ErrUsernameExists:
        return ApiError{
            Code:    ErrorCodeConflict,
            Message: "Resource conflict",
            Detail:  err.Error(),
        }
    case ErrInvalidEmail, ErrInvalidUsername, ErrInvalidPassword:
        return ApiError{
            Code:    ErrorCodeValidation,
            Message: "Validation failed",
            Detail:  err.Error(),
        }
    default:
        return ApiError{
            Code:    ErrorCodeInternal,
            Message: "Internal server error",
            Detail:  "An unexpected error occurred",
        }
    }
}