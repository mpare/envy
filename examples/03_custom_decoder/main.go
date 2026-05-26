package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/mpare/envy"
	"github.com/mpare/envy/decoders"
)

// JSONData is a custom type based on map[string]interface{}.
// It implements the SelfDecoder interface to provide custom JSON parsing.
//
// This allows you to:
// - Parse complex types (JSON, YAML, TOML, etc.) from environment variables
// - Validate complex structures
// - Have full control over parsing logic
//
// Any type can implement SelfDecoder - it just needs the Decode method.
type JSONData map[string]interface{}

// Decode implements the SelfDecoder interface.
// This method is called when envy encounters a JSONData field
// and no built-in decoder matches its type.
func (j *JSONData) Decode(field reflect.Value, raw string, tag decoders.TagReader) error {
	// Parse the raw string as JSON
	if err := json.Unmarshal([]byte(raw), j); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Update the field with the decoded value
	field.Set(reflect.ValueOf(*j))
	return nil
}

// Config demonstrates using custom decoders for complex types.
type Config struct {
	AppName string `env:"APP_NAME,required"`

	// Custom JSON decoder - parses JSON from environment variable
	// Example: FEATURES='{"auth":true,"logging":true,"caching":false}'
	Features JSONData `env:"FEATURES,required"`

	// Custom decoder with default - if env var not set, uses default JSON
	// Example: METADATA='{"version":"1.0.0","env":"prod"}'
	Metadata JSONData `env:"METADATA,default={}"`
}

func main() {
	// For testing, create a map with sample JSON data
	testEnv := map[string]string{
		"APP_NAME": "api-server",
		"FEATURES": `{
			"authentication": true,
			"logging": true,
			"caching": true,
			"rate_limiting": false
		}`,
		"METADATA": `{
			"version": "2.1.0",
			"environment": "production",
			"deployed": "2024-01-15T10:30:00Z"
		}`,
	}

	var cfg Config

	// LoadFrom allows passing a custom environment map (useful for testing)
	if err := envy.LoadFrom(&cfg, testEnv); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("=== Custom Decoder Example ===\n\n")

	fmt.Printf("Application: %s\n\n", cfg.AppName)

	// Display features
	fmt.Println("=== Enabled Features ===")
	for key, value := range cfg.Features {
		status := "disabled"
		if enabled, ok := value.(bool); ok && enabled {
			status = "enabled"
		}
		fmt.Printf("  %s: %s\n", key, status)
	}

	// Display metadata
	if len(cfg.Metadata) > 0 {
		fmt.Println("\n=== Application Metadata ===")
		for key, value := range cfg.Metadata {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Demonstrate how to use the parsed data
	fmt.Println("\n=== Usage Example ===")
	if auth, exists := cfg.Features["authentication"].(bool); exists && auth {
		fmt.Println("✓ Authentication is enabled - initializing auth module...")
	}

	if version, exists := cfg.Metadata["version"].(string); exists {
		fmt.Printf("✓ Application version %s loaded\n", version)
	}

	fmt.Println("\n=== Custom Decoder Pattern ===")
	fmt.Println("To create custom decoders for other types:")
	fmt.Println("1. Define your type (can be a struct, map, or wrapper)")
	fmt.Println("2. Implement SelfDecoder interface with a Decode() method")
	fmt.Println("3. Use the type in your config struct")
	fmt.Println("4. Envy will automatically use your decoder")
	fmt.Println("\nExamples: JSON, YAML, TOML, Protocol Buffers, custom binary formats")
}
