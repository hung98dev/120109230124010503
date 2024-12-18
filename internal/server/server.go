package server

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "hr-backend/pkg/config"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

type Server struct {
    config *config.Config
    router *gin.Engine
}

func NewServer(config *config.Config, router *gin.Engine) *Server {
    return &Server{
        config: config,
        router: router,
    }
}

func (s *Server) Start() error {
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", s.config.ServerPort),
        Handler:      s.router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Server run context
    serverCtx, serverStopCtx := context.WithCancel(context.Background())

    // Listen for syscall signals for process to interrupt/quit
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    go func() {
        <-sig

        // Shutdown signal with grace period of 30 seconds
        shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

        go func() {
            <-shutdownCtx.Done()
            if shutdownCtx.Err() == context.DeadlineExceeded {
                log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
            }
        }()

        // Trigger graceful shutdown
        err := server.Shutdown(shutdownCtx)
        if err != nil {
            log.Fatal().Err(err).Msg("server shutdown error")
        }
        serverStopCtx()
    }()

    // Run the server
    log.Info().Msgf("Server is running on port %d", s.config.ServerPort)
    err := server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        return err
    }

    // Wait for server context to be stopped
    <-serverCtx.Done()
    return nil
}