// internal/middleware/logger.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Stop timer
        timestamp := time.Now()
        latency := timestamp.Sub(start)

        if raw != "" {
            path = path + "?" + raw
        }

        // Log using zerolog
        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Int("status", c.Writer.Status()).
            Str("ip", c.ClientIP()).
            Dur("latency", latency).
            Str("user-agent", c.Request.UserAgent()).
            Msg("Request processed")
    }
}