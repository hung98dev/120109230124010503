package database

import (
    "context"
    "fmt"
    "time"

    "hr-backend/pkg/config"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog/log"
)

func NewDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
    poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("error parsing database config: %w", err)
    }

    // Optimize pool configuration
    poolConfig.MaxConns = 50
    poolConfig.MinConns = 10
    poolConfig.MaxConnLifetime = 30 * time.Minute
    poolConfig.MaxConnIdleTime = 15 * time.Minute
    poolConfig.HealthCheckPeriod = 30 * time.Second
    
    // Add connection timeout
    poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

    // Add query timeout and other connection settings
    poolConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
        // You can add additional connection checks here
        return true
    }

    poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
        // Set various connection parameters
        commands := []string{
            "SET statement_timeout = '10s'",
            "SET lock_timeout = '5s'",
            "SET idle_in_transaction_session_timeout = '15s'",
            "SET client_encoding = 'UTF8'",
            "SET timezone = 'UTC'",
        }

        for _, cmd := range commands {
            if _, err := conn.Exec(ctx, cmd); err != nil {
                return fmt.Errorf("error executing %s: %w", cmd, err)
            }
        }
        return nil
    }

    ctx := context.Background()
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    // Test connection
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("error pinging database: %w", err)
    }

    // Log connection success
    log.Info().
        Str("host", cfg.DBHost).
        Int("port", cfg.DBPort).
        Str("database", cfg.DBName).
        Int32("max_conns", poolConfig.MaxConns).
        Msg("Connected to database")

    // Start connection monitor
    go monitorDBPool(pool)

    return pool, nil
}

// monitorDBPool periodically logs database pool statistics
func monitorDBPool(pool *pgxpool.Pool) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        stats := pool.Stat()
        log.Debug().
            Int32("total_conns", stats.TotalConns()).
            Int32("acquired_conns", stats.AcquiredConns()).
            Int32("idle_conns", stats.IdleConns()).
            Msg("Database pool stats")
    }
}

// CloseDB closes the database connection pool
func CloseDB(pool *pgxpool.Pool) {
    if pool != nil {
        pool.Close()
        log.Info().Msg("Database connection closed")
    }
}

// ExecuteInTransaction executes the given function within a transaction
func ExecuteInTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            // Rollback transaction on panic
            if rbErr := tx.Rollback(ctx); rbErr != nil {
                log.Error().Err(rbErr).Msg("Error rolling back transaction after panic")
            }
            panic(p) // Re-throw panic after rollback
        }
    }()

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            log.Error().Err(rbErr).Msg("Error rolling back transaction")
            return fmt.Errorf("error executing transaction and rolling back: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("error committing transaction: %w", err)
    }

    return nil
}

// HealthCheck performs a health check on the database
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    if err := pool.Ping(ctx); err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }

    var result int
    if err := pool.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
        return fmt.Errorf("database query check failed: %w", err)
    }

    return nil
}

// Example usage of transaction function
func ExampleTransaction(ctx context.Context, pool *pgxpool.Pool) error {
    return ExecuteInTransaction(ctx, pool, func(tx pgx.Tx) error {
        // Perform database operations within transaction
        _, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "John")
        if err != nil {
            return err
        }

        _, err = tx.Exec(ctx, "UPDATE user_counts SET count = count + 1")
        if err != nil {
            return err
        }

        return nil
    })
}