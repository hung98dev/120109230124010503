// pkg/config/config.go
package config

import (
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    // Server
    Environment string
    ServerPort  int

    // Database
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string

    // JWT
    JWTSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    // Load .env file if it exists
    godotenv.Load()

    config := &Config{}

    // Server config
    config.Environment = getEnv("ENV", "development")
    config.ServerPort = getEnvAsInt("PORT", 8080)

    // Database config
    config.DBHost = getEnv("DB_HOST", "localhost")
    config.DBPort = getEnvAsInt("DB_PORT", 5432)
    config.DBUser = getEnv("DB_USER", "postgres")
    config.DBPassword = getEnv("DB_PASSWORD", "545724")
    config.DBName = getEnv("DB_NAME", "hrdb")

    // JWT config
    config.JWTSecret = getEnv("JWT_SECRET", "")
    if config.JWTSecret == "" {
        return nil, fmt.Errorf("JWT_SECRET is required")
    }

    return config, nil
}

// IsProduction checks if the environment is production
func (c *Config) IsProduction() bool {
    return c.Environment == "production"
}

// IsDevelopment checks if the environment is development
func (c *Config) IsDevelopment() bool {
    return c.Environment == "development"
}

// IsTest checks if the environment is test
func (c *Config) IsTest() bool {
    return c.Environment == "test"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        c.DBHost,
        c.DBPort,
        c.DBUser,
        c.DBPassword,
        c.DBName,
    )
}

// Helper function to get an environment variable with a default value
func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultVal
}

// Helper function to get an environment variable as integer with a default value
func getEnvAsInt(key string, defaultVal int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultVal
}

// Helper function to get an environment variable as boolean with a default value
func getEnvAsBool(key string, defaultVal bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultVal
}