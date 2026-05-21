package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mpare/envy"
)

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,default=5432"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Database string `env:"NAME,default=postgres"`
}

// CacheConfig holds cache/Redis settings.
type CacheConfig struct {
	Host    string        `env:"HOST,default=localhost"`
	Port    int           `env:"PORT,default=6379"`
	TTL     time.Duration `env:"TTL,default=1h"`
	Enabled bool          `env:"ENABLED,default=true"`
}

// Config represents the complete microservice configuration.
type Config struct {
	AppName  string        `env:"APP_NAME,required"`
	Port     int           `env:"PORT,default=8080"`
	Debug    bool          `env:"DEBUG,default=false"`
	LogLevel string        `env:"LOG_LEVEL,default=info"`
	Timeout  time.Duration `env:"TIMEOUT,default=30s"`

	// Nested configurations with prefixes
	Database DatabaseConfig `env:",prefix=DB_"`
	Cache    CacheConfig    `env:",prefix=CACHE_"`
}

func main() {
	var cfg Config

	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("=== Microservice Configuration ===\n")
	fmt.Printf("App: %s (debug=%v, log_level=%s)\n", cfg.AppName, cfg.Debug, cfg.LogLevel)
	fmt.Printf("Port: %d, Timeout: %v\n\n", cfg.Port, cfg.Timeout)

	fmt.Printf("Database: %s@%s:%d/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	fmt.Printf("Cache: %s:%d (enabled=%v, ttl=%v)\n", cfg.Cache.Host, cfg.Cache.Port, cfg.Cache.Enabled, cfg.Cache.TTL)
}
