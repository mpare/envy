package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/mpare/envy"
)

// AdvancedConfig demonstrates advanced envy features:
// - Variable expansion with ${VAR}
// - URL parsing
// - Map parsing with custom separators
// - Empty value validation
type AdvancedConfig struct {
	// URL parsing - automatically converts string to url.URL type
	ServerURL url.URL `env:"SERVER_URL,default=http://localhost:8080"`

	// Variable expansion - expands ${VAR} references in values
	// If BASE_PATH=/app and CONFIG_PATH=${BASE_PATH}/config,
	// then ConfigPath will be "/app/config"
	BasePath   string `env:"BASE_PATH,default=/app"`
	ConfigPath string `env:"CONFIG_PATH,default=${BASE_PATH}/config,expand"`

	// notEmpty validation - field must not be set to empty string
	APIKey string `env:"API_KEY,notEmpty"`

	// Map parsing with default separator (key:value pairs separated by comma)
	// Example: "db:5432,cache:6379,api:3000"
	ServicePorts map[string]int `env:"SERVICE_PORTS"`

	// Map with custom separators
	// Example: "Authorization=Bearer token;Content-Type=application/json"
	Headers map[string]string `env:"HEADERS,separator=;,keyValSeparator==="`

	// Slices work as before
	AllowedIPs []string `env:"ALLOWED_IPS,separator=;"`
}

func main() {
	var cfg AdvancedConfig

	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("=== Advanced Configuration Features ===\n")

	// URL parsing
	fmt.Printf("Server URL: %s\n", cfg.ServerURL.String())
	fmt.Printf("  Scheme: %s\n", cfg.ServerURL.Scheme)
	fmt.Printf("  Host: %s\n", cfg.ServerURL.Host)

	// Variable expansion
	fmt.Printf("\nPath Configuration:\n")
	fmt.Printf("  Base Path: %s\n", cfg.BasePath)
	fmt.Printf("  Config Path: %s (expanded from ${BASE_PATH}/config)\n", cfg.ConfigPath)

	// notEmpty validation
	fmt.Printf("\nAPI Key: %s (validated as non-empty)\n", cfg.APIKey)

	// Map parsing
	if len(cfg.ServicePorts) > 0 {
		fmt.Println("\nService Ports (map[string]int):")
		for service, port := range cfg.ServicePorts {
			fmt.Printf("  %s: %d\n", service, port)
		}
	}

	if len(cfg.Headers) > 0 {
		fmt.Println("\nHTTP Headers (map[string]string):")
		for key, value := range cfg.Headers {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	if len(cfg.AllowedIPs) > 0 {
		fmt.Println("\nAllowed IPs:")
		for _, ip := range cfg.AllowedIPs {
			fmt.Printf("  - %s\n", ip)
		}
	}
}
