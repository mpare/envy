package main

import (
	"fmt"
	"log"

	"github.com/mpare/envy"
)

// Config demonstrates basic envy usage with common configuration fields.
type Config struct {
	AppName string `env:"APP_NAME,default=myapp"`
	Port    int    `env:"PORT,default=8080"`
	Debug   bool   `env:"DEBUG,default=false"`
}

func main() {
	var cfg Config

	// Load from environment variables (os.Environ)
	if err := envy.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("Starting %s on port %d (debug=%v)\n", cfg.AppName, cfg.Port, cfg.Debug)
}
