// pkg/utils/jwt.go
package utils

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// CustomClaims định nghĩa cấu trúc của JWT claims
type CustomClaims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerateToken tạo JWT token với user_id
func GenerateToken(userID int, secretKey string) (string, error) {
    // Tạo claims
    claims := CustomClaims{
        userID,
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // Token hết hạn sau 24h
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    // Tạo token với claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Ký token với secret key
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }

    return tokenString, nil
}

// ValidateToken kiểm tra và parse JWT token
func ValidateToken(tokenString string, secretKey string) (*CustomClaims, error) {
    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Kiểm tra signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    // Kiểm tra token có hợp lệ không
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    // Lấy claims từ token
    claims, ok := token.Claims.(*CustomClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }

    return claims, nil
}

// ExtractUserID lấy user_id từ token
func ExtractUserID(tokenString string, secretKey string) (int, error) {
    claims, err := ValidateToken(tokenString, secretKey)
    if err != nil {
        return 0, err
    }
    return claims.UserID, nil
}