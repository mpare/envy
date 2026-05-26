package main

import (
	"fmt"
	"log"

	"github.com/mpare/envy"
)

// Config demonstrates basic envy usage with common configuration fields.
// This is the simplest way to use envy - define a struct with env tags and load it.
//
// Tag options used:
//   - env:"NAME" - the environment variable name to read
//   - default=value - fallback value if env var is not set
type Config struct {
	// AppName reads from APP_NAME, defaults to "myapp" if not set
	AppName string `env:"APP_NAME,default=myapp"`

	// Port reads from PORT, defaults to 8080 if not set
	// Automatically converts string to int
	Port int `env:"PORT,default=8080"`

	// Debug reads from DEBUG, defaults to false
	// Accepts: true, false, 1, 0, t, f (case-insensitive)
	Debug bool `env:"DEBUG,default=false"`
}

func main() {
	var cfg Config

	// Load from os.Environ() - the actual environment variables
	// If any required field is missing or type conversion fails,
	// Load returns an error that can be checked
	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// All fields are now populated with type-safe values
	fmt.Printf("Application: %s\n", cfg.AppName)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug: %v\n", cfg.Debug)

	// For testing, you can use LoadFrom with a custom map:
	// testEnv := map[string]string{
	//     "APP_NAME": "testapp",
	//     "PORT": "9000",
	// }
	// envy.LoadFrom(&cfg, testEnv)
}
