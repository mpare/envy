package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mpare/envy"
)

// DatabaseConfig holds database connection settings.
// This is a nested struct - it gets populated with a prefix (DB_)
// so it looks for DB_HOST, DB_PORT, DB_USER, etc. in the environment.
type DatabaseConfig struct {
	// required - this field must be set or envy returns a validation error
	Host string `env:"HOST,required"`

	// default value - uses 5432 if DB_PORT is not set
	Port int `env:"PORT,default=5432"`

	// required - must be provided
	User string `env:"USER,required"`

	// required - often kept in env for security
	Password string `env:"PASSWORD,required"`

	// default value
	Database string `env:"NAME,default=postgres"`
}

// CacheConfig holds cache/Redis settings.
// Also a nested struct with CACHE_ prefix.
type CacheConfig struct {
	Host    string        `env:"HOST,default=localhost"`
	Port    int           `env:"PORT,default=6379"`
	TTL     time.Duration `env:"TTL,default=1h"` // time.Duration type
	Enabled bool          `env:"ENABLED,default=true"`
}

// Config represents the complete microservice configuration.
// Demonstrates combining multiple nested configs with different prefixes.
type Config struct {
	// Top-level fields
	AppName string `env:"APP_NAME,required"`
	Port    int    `env:"PORT,default=8080"`
	Debug   bool   `env:"DEBUG,default=false"`

	// String with choices (you handle the validation in your code)
	LogLevel string `env:"LOG_LEVEL,default=info"`

	// Duration type - automatically parsed from strings like "30s", "5m", "2h"
	Timeout time.Duration `env:"TIMEOUT,default=30s"`

	// Nested struct with prefix - looks for DB_* environment variables
	// This allows organizing configuration hierarchically
	Database DatabaseConfig `env:",prefix=DB_"`

	// Another nested struct with different prefix - looks for CACHE_* environment variables
	Cache CacheConfig `env:",prefix=CACHE_"`
}

func main() {
	var cfg Config

	// Load will search for all environment variables with appropriate prefixes
	// Example environment variables:
	//   APP_NAME=myservice
	//   PORT=9000
	//   DEBUG=true
	//   LOG_LEVEL=debug
	//   TIMEOUT=1m
	//   DB_HOST=localhost
	//   DB_USER=admin
	//   DB_PASSWORD=secret
	//   CACHE_HOST=redis.local
	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("=== Microservice Configuration ===")
	fmt.Println()

	// Application settings
	fmt.Printf("Application: %s\n", cfg.AppName)
	fmt.Printf("  Port: %d\n", cfg.Port)
	fmt.Printf("  Debug: %v\n", cfg.Debug)
	fmt.Printf("  Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("  Request Timeout: %v\n\n", cfg.Timeout)

	// Database settings
	fmt.Printf("Database Configuration:\n")
	fmt.Printf("  Connection: %s@%s:%d/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	fmt.Printf("  Password: %s\n\n", cfg.Database.Password)

	// Cache settings
	fmt.Printf("Cache Configuration:\n")
	fmt.Printf("  Host: %s:%d\n", cfg.Cache.Host, cfg.Cache.Port)
	fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
	fmt.Printf("  TTL: %v\n", cfg.Cache.TTL)

	// In a real application, you would now use these config values:
	// - Connect to database using cfg.Database
	// - Connect to cache using cfg.Cache
	// - Set up logging with cfg.LogLevel
	// - Start HTTP server on cfg.Port
}
