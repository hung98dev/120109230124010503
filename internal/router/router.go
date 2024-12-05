package router

import (
    "time"

    "hr-backend/internal/auth/handler"
    "hr-backend/internal/middleware"
    "hr-backend/pkg/config"
    "github.com/gin-gonic/gin"
)

type Router struct {
    config      *config.Config
    authHandler *handler.AuthHandler
}

func NewRouter(
    config *config.Config,
    authHandler *handler.AuthHandler,
) *Router {
    return &Router{
        config:      config,
        authHandler: authHandler,
    }
}

func (r *Router) Setup(engine *gin.Engine) {
    // Recovery middleware
    engine.Use(gin.Recovery())

    // Logger middleware
    engine.Use(middleware.Logger())

    // CORS middleware
    engine.Use(middleware.CORS())

    // Health check
    engine.GET("/health", r.healthCheck)

    // API v1
    v1 := engine.Group("/api/v1")
    {
        // Auth routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", r.authHandler.Register)
            auth.POST("/login", r.authHandler.Login)
        }

        // Protected routes
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(r.config.JWTSecret))
        {
            // Add protected routes here
        }
    }
}

func (r *Router) healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":    "OK",
        "timestamp": time.Now().Format(time.RFC3339),
        "version":   "1.0.0",
        "env":       r.config.Environment,
    })
}