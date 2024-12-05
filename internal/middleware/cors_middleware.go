// internal/middleware/cors_middleware.go
package middleware

import (
    "github.com/gin-gonic/gin"
)

// CORSConfig cấu hình cho CORS middleware
type CORSConfig struct {
    AllowOrigins     []string
    AllowMethods     []string
    AllowHeaders     []string
    ExposeHeaders    []string
    AllowCredentials bool
    MaxAge           string
}

// DefaultCORSConfig trả về cấu hình CORS mặc định
func DefaultCORSConfig() *CORSConfig {
    return &CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{
            "GET",
            "POST",
            "PUT",
            "PATCH",
            "DELETE",
            "HEAD",
            "OPTIONS",
        },
        AllowHeaders: []string{
            "Origin",
            "Content-Length",
            "Content-Type",
            "Authorization",
            "X-Requested-With",
            "Accept",
        },
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           "86400", // 24 hours
    }
}

// CORS middleware handler
func CORS() gin.HandlerFunc {
    config := DefaultCORSConfig()

    return func(c *gin.Context) {
        // Set CORS headers
        origin := c.Request.Header.Get("Origin")
        if origin != "" {
            for _, allowOrigin := range config.AllowOrigins {
                if allowOrigin == "*" || allowOrigin == origin {
                    c.Header("Access-Control-Allow-Origin", origin)
                    break
                }
            }
        }

        // Set other headers
        c.Header("Access-Control-Allow-Methods", join(config.AllowMethods))
        c.Header("Access-Control-Allow-Headers", join(config.AllowHeaders))
        c.Header("Access-Control-Expose-Headers", join(config.ExposeHeaders))
        c.Header("Access-Control-Max-Age", config.MaxAge)

        if config.AllowCredentials {
            c.Header("Access-Control-Allow-Credentials", "true")
        }

        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

// join converts string slice to comma separated string
func join(s []string) string {
    if len(s) == 0 {
        return ""
    }

    result := s[0]
    for i := 1; i < len(s); i++ {
        result += ", " + s[i]
    }
    return result
}