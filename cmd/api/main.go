package main

import (
    "os"

    "hr-backend/internal/auth/handler"
    "hr-backend/internal/auth/repository"
    "hr-backend/internal/auth/service"
    "hr-backend/internal/router"
    "hr-backend/internal/server"
    "hr-backend/pkg/config"
    "hr-backend/pkg/database"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Setup logger
    setupLogger()

    // Load config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal().Err(err).Msg("cannot load config")
    }

    // Set Gin mode
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Connect to database
    db, err := database.NewDBPool(cfg)
    if err != nil {
        log.Fatal().Err(err).Msg("cannot connect to database")
    }
    defer db.Close()

    // Initialize repositories
    authRepo := repository.NewAuthRepository(db)

    // Initialize services
    authService := service.NewAuthService(authRepo, cfg.JWTSecret)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)

    // Initialize router
    routerHandler := router.NewRouter(cfg, authHandler)

    // Create Gin engine
    engine := gin.New()

    // Setup routes
    routerHandler.Setup(engine)

    // Initialize and start server
    srv := server.NewServer(cfg, engine)
    if err := srv.Start(); err != nil {
        log.Fatal().Err(err).Msg("cannot start server")
        os.Exit(1)
    }
}

func setupLogger() {
    output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05"}
    log.Logger = log.Output(output).With().Caller().Logger()

    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if os.Getenv("ENV") == "development" {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
}