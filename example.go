package envy

import (
	"fmt"
	"time"
)

// ExampleUsage demonstrates the typical usage of envy
func ExampleUsage() {
	// Example 1: Basic usage
	type Config struct {
		AppName string        `env:"APP_NAME,default=myapp"`
		Port    int           `env:"PORT,default=8080"`
		Debug   bool          `env:"DEBUG,default=false"`
		Timeout time.Duration `env:"TIMEOUT,default=30s"`
	}

	// In real code, you'd use:
	// var cfg Config
	// envy.MustLoad(&cfg)

	// For this example, we'll use LoadFrom with test data
	var cfg Config
	err := LoadFrom(&cfg, map[string]string{
		"PORT":  "9090",
		"DEBUG": "true",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("App: %s, Port: %d, Debug: %v, Timeout: %v\n",
		cfg.AppName, cfg.Port, cfg.Debug, cfg.Timeout)
	// Output: App: myapp, Port: 9090, Debug: true, Timeout: 30s
}

// ExampleNested demonstrates nested struct usage
func ExampleNested() {
	type DatabaseConfig struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,default=5432"`
		Password string `env:"PASSWORD,required"`
	}

	type Config struct {
		AppName string         `env:"APP_NAME,default=myapp"`
		DB      DatabaseConfig `env:",prefix=DB_"`
	}

	var cfg Config
	err := LoadFrom(&cfg, map[string]string{
		"DB_HOST":     "db.example.com",
		"DB_PASSWORD": "secret123",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("App: %s, DB Host: %s, DB Port: %d\n",
		cfg.AppName, cfg.DB.Host, cfg.DB.Port)
	// Output: App: myapp, DB Host: db.example.com, DB Port: 5432
}

// ExampleSlices demonstrates slice support
func ExampleSlices() {
	type Config struct {
		AllowedIPs []string `env:"ALLOWED_IPS,separator=;"`
		Ports      []int    `env:"PORTS,separator=,"`
	}

	var cfg Config
	err := LoadFrom(&cfg, map[string]string{
		"ALLOWED_IPS": "192.168.1.1;192.168.1.2;192.168.1.3",
		"PORTS":       "8080,8081,8082",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("IPs: %v, Ports: %v\n", cfg.AllowedIPs, cfg.Ports)
	// Output: IPs: [192.168.1.1 192.168.1.2 192.168.1.3], Ports: [8080 8081 8082]
}

// ExampleErrorHandling demonstrates error handling
func ExampleErrorHandling() {
	type Config struct {
		Secret string `env:"SECRET,required"`
		Port   int    `env:"PORT"`
	}

	var cfg Config
	err := LoadFrom(&cfg, map[string]string{
		"PORT": "invalid",
	})

	if err != nil {
		if ve, ok := err.(*ValidationError); ok {
			for _, fieldErr := range ve.Errors {
				fmt.Printf("%s (%s): %s\n",
					fieldErr.Field, fieldErr.EnvKey, fieldErr.Message)
			}
		}
	}
	// Output:
	// Secret (SECRET): required but not set
	// Port (PORT): invalid int "invalid": strconv.ParseInt: parsing "invalid": invalid syntax
}

// ExampleTypes demonstrates various supported types
func ExampleTypes() {
	type Config struct {
		AppName    string        `env:"APP_NAME,default=myapp"`
		Port       int           `env:"PORT,default=8080"`
		Rate       float64       `env:"RATE,default=1.0"`
		Debug      bool          `env:"DEBUG,default=false"`
		Timeout    time.Duration `env:"TIMEOUT,default=30s"`
		Tags       []string      `env:"TAGS,separator=;"`
		Thresholds []float64     `env:"THRESHOLDS,separator=,"`
	}

	var cfg Config
	err := LoadFrom(&cfg, map[string]string{
		"RATE":       "2.5",
		"DEBUG":      "true",
		"TIMEOUT":    "1m",
		"TAGS":       "prod;monitoring;critical",
		"THRESHOLDS": "0.5,0.75,0.95",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("App: %s\n", cfg.AppName)
	fmt.Printf("Port: %d, Rate: %f, Debug: %v\n", cfg.Port, cfg.Rate, cfg.Debug)
	fmt.Printf("Timeout: %v\n", cfg.Timeout)
	fmt.Printf("Tags: %v\n", cfg.Tags)
	fmt.Printf("Thresholds: %v\n", cfg.Thresholds)
	// Output:
	// App: myapp
	// Port: 8080, Rate: 2.500000, Debug: true
	// Timeout: 1m0s
	// Tags: [prod monitoring critical]
	// Thresholds: [0.5 0.75 0.95]
}
