package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/mpare/envy"
	"github.com/mpare/envy/decoders"
)

// JSONData is a custom type that implements the SelfDecoder interface.
// It allows parsing JSON from environment variables.
type JSONData map[string]interface{}

// Decode implements the SelfDecoder interface for custom JSON parsing.
func (j *JSONData) Decode(field reflect.Value, raw string, tag decoders.TagReader) error {
	if err := json.Unmarshal([]byte(raw), j); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	field.Set(reflect.ValueOf(*j))
	return nil
}

// Config demonstrates using custom decoders for complex types.
type Config struct {
	AppName  string   `env:"APP_NAME,required"`
	Features JSONData `env:"FEATURES,required"`
	Metadata JSONData `env:"METADATA,default={}"`
}

func main() {
	var cfg Config

	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("Application: %s\n\n", cfg.AppName)

	fmt.Println("=== Features ===")
	for key, value := range cfg.Features {
		fmt.Printf("  %s: %v\n", key, value)
	}

	if len(cfg.Metadata) > 0 {
		fmt.Println("\n=== Metadata ===")
		for key, value := range cfg.Metadata {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}
