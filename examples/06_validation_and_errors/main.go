package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mpare/envy"
)

// ValidationConfig demonstrates various validation features:
// - required: field must be set
// - notEmpty: field must not be empty
// - default: fallback value if not set
type ValidationConfig struct {
	// required - must be set or get a validation error
	AppName string `env:"APP_NAME,required"`

	// notEmpty - if set, must not be empty string
	APIKey string `env:"API_KEY,notEmpty"`

	// optional with default
	Port int `env:"PORT,default=8080"`

	// required with no default
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Complex types with defaults
	Timeout time.Duration `env:"TIMEOUT,default=30s"`
	Retries int           `env:"RETRIES,default=3"`
	Debug   bool          `env:"DEBUG,default=false"`
}

func main() {
	// Example 1: With valid configuration
	fmt.Println("=== Example 1: Valid Configuration ===\n")
	validEnv := map[string]string{
		"APP_NAME":     "myapp",
		"API_KEY":      "secret-key-123",
		"DATABASE_URL": "postgres://localhost/mydb",
	}

	var cfg1 ValidationConfig
	if err := envy.LoadFrom(&cfg1, validEnv); err != nil {
		log.Fatalf("valid config failed: %v", err)
	}

	fmt.Printf("✓ Configuration loaded successfully:\n")
	fmt.Printf("  AppName: %s\n", cfg1.AppName)
	fmt.Printf("  APIKey: %s\n", cfg1.APIKey)
	fmt.Printf("  Port: %d\n", cfg1.Port)
	fmt.Printf("  DatabaseURL: %s\n", cfg1.DatabaseURL)
	fmt.Printf("  Timeout: %v\n", cfg1.Timeout)
	fmt.Printf("  Retries: %d\n", cfg1.Retries)

	// Example 2: Missing required field
	fmt.Println("\n=== Example 2: Missing Required Field ===\n")
	invalidEnv := map[string]string{
		"APP_NAME": "myapp",
		// DATABASE_URL is missing!
	}

	var cfg2 ValidationConfig
	if err := envy.LoadFrom(&cfg2, invalidEnv); err != nil {
		var ve *envy.ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("✗ Validation errors found (%d errors):\n", len(ve.Errors))
			for _, fieldErr := range ve.Errors {
				fmt.Printf("  - Field: %s\n", fieldErr.Field)
				fmt.Printf("    EnvKey: %s\n", fieldErr.EnvKey)
				fmt.Printf("    Message: %s\n\n", fieldErr.Message)
			}
		}
	}

	// Example 3: Empty value with notEmpty validation
	fmt.Println("=== Example 3: Empty Value with notEmpty ===\n")
	emptyEnv := map[string]string{
		"APP_NAME":     "myapp",
		"API_KEY":      "", // Empty! This should fail notEmpty validation
		"DATABASE_URL": "postgres://localhost/mydb",
	}

	var cfg3 ValidationConfig
	if err := envy.LoadFrom(&cfg3, emptyEnv); err != nil {
		var ve *envy.ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("✗ Validation errors found (%d errors):\n", len(ve.Errors))
			for _, fieldErr := range ve.Errors {
				fmt.Printf("  - %s (%s): %s\n", fieldErr.Field, fieldErr.EnvKey, fieldErr.Message)
			}
		}
	}

	// Example 4: Multiple validation errors
	fmt.Println("=== Example 4: Multiple Validation Errors ===\n")
	multiErrorEnv := map[string]string{
		// APP_NAME missing (required)
		"API_KEY": "", // Empty (notEmpty)
		// DATABASE_URL missing (required)
	}

	var cfg4 ValidationConfig
	if err := envy.LoadFrom(&cfg4, multiErrorEnv); err != nil {
		var ve *envy.ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("✗ Multiple validation errors found (%d errors):\n\n", len(ve.Errors))
			for i, fieldErr := range ve.Errors {
				fmt.Printf("  Error %d:\n", i+1)
				fmt.Printf("    Field: %s\n", fieldErr.Field)
				fmt.Printf("    EnvKey: %s\n", fieldErr.EnvKey)
				fmt.Printf("    Message: %s\n\n", fieldErr.Message)
			}
		}
	}

	fmt.Println("=== Validation Features Summary ===\n")
	fmt.Println("Supported validations:")
	fmt.Println("  - required: field must be set (error if missing)")
	fmt.Println("  - notEmpty: field must not be empty string (error if empty)")
	fmt.Println("  - default=value: fallback value if not set")
	fmt.Println("\nAll validation errors are collected and reported together,")
	fmt.Println("so you see all problems at once instead of failing on first error.")
}
