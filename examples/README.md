# Envy Examples

This directory contains runnable examples demonstrating different ways to use envy.

Each example is self-contained and includes detailed comments explaining the features.

## Examples Overview

### 1. [Basic Configuration](01_basic)

**Use this when:** Getting started with envy for simple applications

Demonstrates:
- Loading basic types (string, int, bool)
- Using default values
- Simple type conversion
- Testing with custom environment maps

**Run:** `cd 01_basic && go run main.go`

Environment variables:
```bash
export APP_NAME=myapp
export PORT=8080
export DEBUG=false
```

---

### 2. [Microservice Configuration](02_microservice)

**Use this when:** Building microservices with multiple configuration sections

Demonstrates:
- Nested configuration structs
- Environment variable prefixes (`prefix=`)
- Required vs optional fields
- Duration parsing (`time.Duration`)
- Organizing related settings

**Run:** `cd 02_microservice && go run main.go`

Environment variables:
```bash
export APP_NAME=my-service
export PORT=9000
export DB_HOST=localhost
export DB_USER=admin
export DB_PASSWORD=secret
export CACHE_HOST=redis.local
export CACHE_TTL=1h
```

---

### 3. [Custom Decoders](03_custom_decoder)

**Use this when:** You need to parse complex types like JSON or YAML

Demonstrates:
- Implementing the `SelfDecoder` interface
- Parsing JSON from environment variables
- Custom parsing logic
- Error handling for invalid data
- Complex type validation

**Run:** `cd 03_custom_decoder && go run main.go`

Environment variables:
```bash
export APP_NAME=api-server
export FEATURES='{"auth":true,"logging":true,"caching":false}'
export METADATA='{"version":"1.0.0","env":"prod"}'
```

---

### 4. [Advanced Features](04_advanced_features)

**Use this when:** Using new advanced features like variable expansion, maps, and URLs

Demonstrates:
- **Variable expansion** (`expand` tag): `${VAR}` references in values
- **URL parsing** (`url.URL` type): Automatic URL parsing
- **Map parsing** (`map[K]V` types): Key-value pairs from env vars
- **Custom separators**: Different delimiters for maps and slices
- **Slice parsing** with custom separators

**Run:** `cd 04_advanced_features && go run main.go`

Environment variables:
```bash
export SERVER_URL=https://api.example.com:8443/v1
export BASE_PATH=/app
export CONFIG_PATH='${BASE_PATH}/config'
export API_KEY=your-secret-key
export SERVICE_PORTS="db:5432,cache:6379,api:3000"
export HEADERS="Authorization=Bearer token;Content-Type=application/json"
export ALLOWED_IPS="192.168.1.0;10.0.0.0;127.0.0.1"
```

---

### 5. [File-Based Configuration](05_file_based_config)

**Use this when:** Loading configuration from files (certificates, configs, etc.)

Demonstrates:
- **File reading** (`file` tag): Load content from file paths
- Combining file reading with defaults
- Multi-line configuration
- Secrets management (certificates, keys)

**Run:** `cd 05_file_based_config && go run main.go`

Environment variables:
```bash
export APP_NAME=myapp
export TLS_CERT=/etc/ssl/certs/server.crt
export PRIVATE_KEY=/etc/ssl/private/server.key
export DB_CONFIG_FILE=/etc/config/database.conf
```

---

### 6. [Validation and Error Handling](06_validation_and_errors)

**Use this when:** Understanding validation rules and error handling

Demonstrates:
- **Required fields** (`required` tag): Fields that must be set
- **Empty validation** (`notEmpty` tag): Fields that must not be empty
- **Default values** (`default=` tag): Fallback values
- **Aggregated errors**: All validation errors collected at once
- **Error inspection**: Examining field-level errors

**Run:** `cd 06_validation_and_errors && go run main.go`

The example shows:
- Valid configuration loading
- Missing required field errors
- Empty value validation errors
- Multiple validation errors at once

---

## Feature Matrix

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
