# Envy Examples

This directory contains runnable examples demonstrating different ways to use envy.

## Examples

### 1. [Basic Configuration](01_basic)

**Use this when:** Getting started with envy for simple applications

Demonstrates:
- Loading basic types (string, int, bool)
- Using default values
- Simple type conversion

**Run:** `cd 01_basic && go run main.go`

---

### 2. [Microservice Configuration](02_microservice)

**Use this when:** Building microservices with multiple configuration sections

Demonstrates:
- Nested configuration structs
- Environment variable prefixes
- Required vs optional fields
- Duration parsing
- Organizing related settings

**Run:** `cd 02_microservice && go run main.go`

---

### 3. [Custom Decoders](03_custom_decoder)

**Use this when:** You need to parse complex types like JSON or YAML

Demonstrates:
- Implementing the `SelfDecoder` interface
- Parsing JSON from environment variables
- Custom parsing logic
- Error handling for invalid data

**Run:** `cd 03_custom_decoder && go run main.go`

---

## Quick Start

### 1. Hello World

```bash
cd 01_basic
PORT=3000 DEBUG=true go run main.go
```

### 2. Full Microservice

```bash
cd 02_microservice
DB_HOST=localhost DB_USER=admin DB_PASSWORD=pass APP_NAME=myapp go run main.go
```

### 3. JSON Configuration

```bash
cd 03_custom_decoder
APP_NAME=myapp FEATURES='{"auth":true,"api":true}' go run main.go
```

## Learning Path

1. **Start with Basic** - Understand simple types and defaults
2. **Move to Microservice** - Learn about nested configs and prefixes
3. **Explore Custom Decoders** - Implement custom parsing for complex types

## Common Patterns

### Environment Variable Naming

```
[COMPONENT_]SETTING_NAME

Examples:
  APP_NAME          # Top-level settings (no prefix)
  DB_HOST           # Database prefix
  CACHE_TTL         # Cache prefix
  LOG_LEVEL         # Service-wide settings
```

### Configuration Organization

```go
type AppConfig struct {
  // Application settings
  Name string `env:"APP_NAME"`
  Port int    `env:"PORT"`
  
  // Sub-components with prefixes
  DB    DatabaseConfig `env:",prefix=DB_"`
  Cache CacheConfig    `env:",prefix=CACHE_"`
}
```

### Testing with LoadFrom

```go
var cfg Config
err := envy.LoadFrom(&cfg, map[string]string{
  "PORT": "8080",
  "DB_HOST": "localhost",
})
```
