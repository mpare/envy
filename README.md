# envy — Typed Environment Variables for Go

A lightweight Go library for loading environment variables into typed structs with validation, defaults, and nested support. Zero external dependencies, stdlib only.

[![Go Version](https://img.shields.io/github/go-mod/go-version/mpare/envy)](https://github.com/mpare/envy)
[![CI](https://github.com/mpare/envy/actions/workflows/ci.yml/badge.svg)](https://github.com/mpare/envy/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mpare/envy)](https://goreportcard.com/report/github.com/mpare/envy)
[![License](https://img.shields.io/github/license/mpare/envy)](LICENSE)

## Features

- ✅ **Type-safe**: String, int, float, bool, duration, slices, maps, URLs, and nested structs
- ✅ **Default values**: Specify defaults with `default=value` tag
- ✅ **Required fields**: Mark fields as `required` and get validation errors
- ✅ **Aggregated errors**: Collect all validation errors in one pass
- ✅ **Nested structs**: Use `prefix=` for environment variable namespacing
- ✅ **Custom separators**: Configure slice/map separators with `separator=`
- ✅ **Variable expansion**: Expand environment variables in values with `expand` tag
- ✅ **File reading**: Load configuration from files with `file` tag
- ✅ **Empty validation**: Enforce non-empty values with `notEmpty` tag
- ✅ **Map support**: Parse maps with configurable key-value separators
- ✅ **URL parsing**: Built-in support for `url.URL` type
- ✅ **Zero dependencies**: Uses Go stdlib only
- ✅ **Test-friendly**: `LoadFrom()` accepts custom env maps for unit tests

## Use Cases

- Load 12-factor app config from environment
- Type-safe configuration management
- Environment variable validation
- Go microservices configuration

## Installation

```bash
go get github.com/mpare/envy
```

## Quick Start

```go
package main

import (
	"log"
	"github.com/mpare/envy"
)

type Config struct {
	AppName    string `env:"APP_NAME,default=myapp"`
	Port       int    `env:"PORT,default=8080"`
	Debug      bool   `env:"DEBUG,default=false"`
}

func main() {
	var cfg Config
	envy.MustLoad(&cfg)
	// cfg is ready with all types correctly parsed
	log.Printf("Starting %s on port %d", cfg.AppName, cfg.Port)
}
```

## Examples

For more comprehensive examples, see the [examples](examples) directory:

- [Basic Configuration](examples/01_basic) - Simple types and defaults
- [Microservice Setup](examples/02_microservice) - Nested structs with prefixes
- [Custom Decoders](examples/03_custom_decoder) - Parsing JSON from environment variables

Each example is runnable and documented with detailed explanations.

## Struct Tag Format

```go
type Config struct {
	// Basic string
	AppName string `env:"APP_NAME"`

	// With default value
	Port int `env:"PORT,default=8080"`

	// Required field (error if not set)
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Must not be empty
	Secret string `env:"SECRET,notEmpty"`

	// Expand environment variables in value (e.g., FOO_${BAR})
	Path string `env:"CUSTOM_PATH,expand"`

	// Read value from file path
	CertFile string `env:"CERT_PATH,file"`

	// Various types
	Debug      bool          `env:"DEBUG,default=false"`
	Timeout    time.Duration `env:"TIMEOUT,default=30s"`
	MaxRetries int           `env:"MAX_RETRIES,default=3"`
	Rate       float64       `env:"RATE,default=1.5"`
	ServerURL  url.URL       `env:"SERVER_URL"`

	// Slice with custom separator
	AllowedIPs []string `env:"ALLOWED_IPS,separator=;"`
	Workers    []int    `env:"WORKER_PORTS,separator=,"`

	// Map with custom separators (item separator and key:value separator)
	Labels map[string]string `env:"LABELS,separator=;,keyValSeparator=:"`

	// Nested struct with prefix
	Database DatabaseConfig `env:",prefix=DB_"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,default=5432"`
	Password string `env:"PASSWORD,required"`
}
// This looks for: DB_HOST, DB_PORT, DB_PASSWORD

// Tag Options Reference:
// - env:"KEY"              : Set environment variable name (required)
// - ,required              : Mark field as required (error if not set)
// - ,default=value         : Provide default value if env var not set
// - ,notEmpty              : Error if env var is set but empty
// - ,expand                : Expand ${VAR} references in the value
// - ,file                  : Read value from file path specified in env var
// - ,separator=,           : Custom separator for slices and map items (default: ",")
// - ,keyValSeparator=:     : Custom separator for map key:value pairs (default: ":")
// - ,prefix=PREFIX_        : For nested structs, adds prefix to all env var names
```

**Important**: All struct fields must be **exported** (start with a capital letter). Unexported fields are silently skipped.

## API

### Load(destination any) error

Loads environment variables from `os.Environ()` into the destination struct.

```go
var cfg Config
if err := envy.Load(&cfg); err != nil {
	log.Fatal(err)
}
```

### MustLoad(destination any)

Same as `Load()` but panics on error. Useful at application startup.

```go
var cfg Config
envy.MustLoad(&cfg) // Panics if validation fails
```

### LoadFrom(destination any, environ map[string]string) error

Loads from a custom map. Useful for testing.

```go
var cfg Config
err := envy.LoadFrom(&cfg, map[string]string{
	"PORT":     "9090",
	"DATABASE_URL": "postgres://localhost/test",
})
```

## Supported Types

| Type | Example Value |
|------|---------------|
| `string` | `"hello"` |
| `int`, `int8`, `int16`, `int32`, `int64` | `"8080"` |
| `float32`, `float64` | `"1.5"` |
| `bool` | `"true"`, `"false"`, `"1"`, `"0"`, `"t"`, `"f"` |
| `time.Duration` | `"30s"`, `"5m"`, `"1h30m"` |
| `url.URL` | `"https://example.com:8080/path"` |
| `[]string` | `"a,b,c"` or `"a;b;c"` (custom separator via `separator=`) |
| `[]int` | `"1,2,3"` (custom separator via `separator=`) |
| `[]float32`, `[]float64` | `"1.5,2.5,3.5"` (custom separator via `separator=`) |
| `[]bool` | `"true,false,1"` (custom separator via `separator=`) |
| `map[string]string` | `"key1:val1,key2:val2"` (separator and keyValSeparator customizable) |
| `map[string]int` | `"db:5432,cache:6379"` |
| Nested Struct | Via `prefix=` for namespacing |
| Custom Types | Types implementing `SelfDecoder` |

## Custom Decoders

For types not natively supported, implement the `SelfDecoder` interface to provide custom decoding logic:

```go
package main

import (
	"encoding/json"
	"reflect"

	"github.com/mpare/envy"
	"github.com/mpare/envy/decoders"
)

// JSONData is a custom type that decodes JSON
type JSONData map[string]interface{}

// Decode implements the SelfDecoder interface
func (j *JSONData) Decode(field reflect.Value, raw string, tag decoders.TagReader) error {
	if err := json.Unmarshal([]byte(raw), j); err != nil {
		return err
	}
	field.Set(reflect.ValueOf(*j))
	return nil
}

// Usage
type Config struct {
	Metadata JSONData `env:"METADATA"`
}

func main() {
	var cfg Config
	envy.MustLoad(&cfg)
	// cfg.Metadata now contains parsed JSON data
}
```

When loading from environment variables:

```bash
export METADATA='{"version":"1.0","env":"prod","features":["auth","logging"]}'
```

**Key points:**
- Implement `decoders.SelfDecoder` interface with a single `Decode` method
- Your type is automatically detected and used when no built-in decoder matches
- Custom decoders receive the raw string value and must parse/validate it
- Return an error if the value is invalid; envy will collect it as a `ValidationError`

## Advanced Features

### Variable Expansion

Use the `expand` tag to expand environment variables within values:

```go
type Config struct {
	BasePath string `env:"BASE_PATH"`
	// If BASE_PATH=/data and CUSTOM_PATH=${BASE_PATH}/config
	// then CustomPath will be "/data/config"
	CustomPath string `env:"CUSTOM_PATH,expand"`
}
```

Environment setup:
```bash
export BASE_PATH=/data
export CUSTOM_PATH='${BASE_PATH}/config'
```

### File Reading

Use the `file` tag to read configuration from a file:

```go
type Config struct {
	// CertPath contains the path to a certificate file
	// The environment variable should contain the file path,
	// and the field will contain the file's content
	Cert string `env:"CERT_PATH,file"`
}
```

Environment setup:
```bash
export CERT_PATH=/etc/ssl/certs/server.crt
```

### Empty Value Validation

Use the `notEmpty` tag to require non-empty values:

```go
type Config struct {
	// Error if SECRET is not set OR is set to an empty string
	Secret string `env:"SECRET,notEmpty"`
}
```

### Map Parsing

Parse maps with customizable separators:

```go
type Config struct {
	// Parse: "key1:val1,key2:val2"
	Labels map[string]string `env:"LABELS"`

	// Custom separators: "key1=val1;key2=val2"
	Headers map[string]string `env:"HEADERS,separator=;,keyValSeparator==="`

	// Parse port mappings: "db:5432,cache:6379"
	Ports map[string]int `env:"SERVICE_PORTS"`
}
```

Environment setup:
```bash
export LABELS="env:prod,version:1.0"
export HEADERS="Authorization=Bearer token;Content-Type=application/json"
export SERVICE_PORTS="db:5432,cache:6379"
```

## Field Behavior

- **Exported vs Unexported**: Only exported (capitalized) fields are processed. Unexported fields are silently skipped.
- **Unset fields without defaults**: Fields that are not required and have no default value retain their **zero value** (empty string, 0, false, nil slice, etc.)
- **Required fields**: If marked `required` and not found in the environment, a validation error is collected.
- **Nested structs**: Nested struct fields must also be exported. Each nested struct looks for environment variables prefixed with the `prefix` tag value.

## Error Handling

Validation errors are collected and reported together:

```go
type Config struct {
	Secret   string `env:"SECRET,required"`
	Port     int    `env:"PORT"`
	Timeout  time.Duration `env:"TIMEOUT"`
}

var cfg Config
err := envy.LoadFrom(&cfg, map[string]string{
	"PORT":    "invalid",
	"TIMEOUT": "bad",
})

if err != nil {
	var ve *envy.ValidationError
	if errors.As(err, &ve) {
		for _, fieldErr := range ve.Errors {
			fmt.Printf("%s (%s): %s\n", fieldErr.Field, fieldErr.EnvKey, fieldErr.Message)
		}
	}
}

// Output:
// Secret (SECRET): required but not set
// Port (PORT): invalid int "invalid": strconv.ParseInt: parsing "invalid": invalid syntax
// Timeout (TIMEOUT): invalid duration "bad": time.ParseDuration: ...
```

## Examples

### Basic Configuration

```go
type Config struct {
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,default=5432"`
}

var cfg Config
envy.MustLoad(&cfg)
```

### Database with Nested Struct

```go
type Config struct {
	Database struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,default=5432"`
		User     string `env:"USER,required"`
		Password string `env:"PASSWORD,required"`
	} `env:",prefix=DB_"`
}

// Looks for: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD
var cfg Config
envy.MustLoad(&cfg)
```

### Multiple Environments

```go
// Environment variables:
// APP_PORT=8080
// APP_DEBUG=true
// CACHE_HOST=redis:6379
// CACHE_TTL=1h

type Config struct {
	App struct {
		Port  int  `env:"PORT"`
		Debug bool `env:"DEBUG"`
	} `env:",prefix=APP_"`
	Cache struct {
		Host string        `env:"HOST"`
		TTL  time.Duration `env:"TTL"`
	} `env:",prefix=CACHE_"`
}

var cfg Config
envy.MustLoad(&cfg)
```

### Comma-separated Lists

```go
type Config struct {
	AllowedCIDRs []string `env:"ALLOWED_CIDRS,separator=,"`
	Servers      []string `env:"SERVERS,separator=;"`
	Ports        []int    `env:"PORTS,separator=,"`
}

// ALLOWED_CIDRS="10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
// SERVERS="server1;server2;server3"
// PORTS="8080,8081,8082"

var cfg Config
envy.MustLoad(&cfg)
```

## Tag Reference

| Tag | Format | Example | Description |
|-----|--------|---------|-------------|
| `env` | `"KEY"` | `env:"PORT"` | Environment variable name (required) |
| `default` | `default=value` | `env:"PORT,default=8080"` | Default value if not set (used as literal string, not split for slices) |
| `required` | `required` | `env:"API_KEY,required"` | Must be set; error if missing |
| `separator` | `separator=char` | `env:"IPS,separator=;"` | Separator for slice types (default: `,`) |
| `prefix` | `prefix=PREFIX_` | `env:",prefix=DB_"` | Prefix for nested struct fields |

Multiple options are comma-separated:
```go
Port int `env:"PORT,default=8080"` // Single option
Tags []string `env:"TAGS,separator=;"` // Custom separator
URL string `env:"DATABASE_URL,required"` // Required
```

**Note on combinations**: Using both `default` and `required` is semantically contradictory (if there's a default, the field is never truly required). If both are specified, `default` takes precedence—the field will use the default value if the env var is not set.

## Testing

Use `LoadFrom()` to inject test data:

```go
func TestMyService(t *testing.T) {
	var cfg Config
	err := envy.LoadFrom(&cfg, map[string]string{
		"PORT": "9090",
		"DEBUG": "true",
		"DATABASE_URL": "postgres://localhost/testdb",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test with cfg
	service := NewService(cfg)
	// ...
}
```

## License

MIT

## Contributing

Contributions welcome! Please ensure all tests pass:

```bash
go test ./...
```
