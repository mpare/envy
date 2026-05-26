package main

import (
	"fmt"
	"log"

	"github.com/mpare/envy"
)

// FileBasedConfig demonstrates the file reading feature.
// The `file` tag tells envy to read the environment variable value as a file path
// and load the file content into the field.
//
// This is useful for:
// - Loading certificates and keys from files
// - Reading multi-line configuration from files
// - Keeping secrets in separate files for security
type FileBasedConfig struct {
	// Regular string field
	AppName string `env:"APP_NAME,default=myapp"`

	// File reading: reads content from file path specified in env var
	// Example: TLS_CERT=/etc/ssl/certs/server.crt
	// Then CertData will contain the file contents
	TLSCert string `env:"TLS_CERT,file"`

	// You can combine file with default values
	// If the env var isn't set, the default will be used as-is (not as file path)
	PrivateKey string `env:"PRIVATE_KEY,default=-----BEGIN PRIVATE KEY-----\nMIIEvQ..."`

	// Another file-based config
	DatabaseDSN string `env:"DB_CONFIG_FILE,file,default=/tmp/default-db-config"`
}

func main() {
	var cfg FileBasedConfig

	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("=== File-Based Configuration ===")
	fmt.Println()

	fmt.Printf("Application: %s\n", cfg.AppName)
	fmt.Println()

	if cfg.TLSCert != "" {
		fmt.Printf("TLS Certificate (from file):\n")
		// In a real scenario, you'd see the actual certificate content
		if len(cfg.TLSCert) > 50 {
			fmt.Printf("  %s... (%d bytes)\n\n", cfg.TLSCert[:50], len(cfg.TLSCert))
		} else {
			fmt.Printf("  %s\n\n", cfg.TLSCert)
		}
	}

	if cfg.PrivateKey != "" {
		fmt.Printf("Private Key:\n")
		if len(cfg.PrivateKey) > 50 {
			fmt.Printf("  %s... (%d bytes)\n\n", cfg.PrivateKey[:50], len(cfg.PrivateKey))
		} else {
			fmt.Printf("  %s\n\n", cfg.PrivateKey)
		}
	}

	if cfg.DatabaseDSN != "" {
		fmt.Printf("Database Configuration (from file):\n")
		if len(cfg.DatabaseDSN) > 80 {
			fmt.Printf("  %s... (%d bytes)\n", cfg.DatabaseDSN[:80], len(cfg.DatabaseDSN))
		} else {
			fmt.Printf("  %s\n", cfg.DatabaseDSN)
		}
	}
}
