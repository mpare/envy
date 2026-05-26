# Advanced Features Example

This example demonstrates advanced envy features added in recent versions.

## Features Covered

### 1. Variable Expansion (`expand` tag)

Expand `${VAR_NAME}` references in environment variable values.

```go
BasePath   string `env:"BASE_PATH,default=/app"`
ConfigPath string `env:"CONFIG_PATH,default=${BASE_PATH}/config,expand"`
```

**Environment setup:**
```bash
export BASE_PATH=/data
export CONFIG_PATH='${BASE_PATH}/config'
# Result: ConfigPath = "/data/config"
```

**Use cases:**
- Composing paths based on other variables
- Building URLs with dynamic base addresses
- Referencing other environment variables for consistency

---

### 2. URL Parsing (`url.URL` type)

Automatically parse URLs into `net/url.URL` structs.

```go
ServerURL url.URL `env:"SERVER_URL,default=http://localhost:8080"`
```

**Environment setup:**
```bash
export SERVER_URL=https://api.example.com:8443/v1
```

**Access parsed components:**
```go
fmt.Println(cfg.ServerURL.Scheme)   // "https"
fmt.Println(cfg.ServerURL.Host)     // "api.example.com:8443"
fmt.Println(cfg.ServerURL.Path)     // "/v1"
```

---

### 3. Map Parsing

Parse maps with configurable separators.

**Default separators:** item separator `,` and key-value separator `:`

```go
ServicePorts map[string]int `env:"SERVICE_PORTS"`
```

**Environment setup:**
```bash
export SERVICE_PORTS="db:5432,cache:6379,api:3000"
```

**Custom separators:**

```go
Headers map[string]string `env:"HEADERS,separator=;,keyValSeparator==="`
```

**Environment setup:**
```bash
export HEADERS="Authorization=Bearer token;Content-Type=application/json"
```

**Use cases:**
- Service discovery (service name to port mapping)
- HTTP headers and metadata
- Feature flags and configuration options
- Dynamic port allocation

---

### 4. Empty Value Validation (`notEmpty` tag)

Validate that fields are not empty strings.

```go
APIKey string `env:"API_KEY,notEmpty"`
```

**Difference from `required`:**
- `required`: Field must be set (error if missing)
- `notEmpty`: Field must not be set to empty string (error if `=""`)

**Environment setup:**
```bash
# Valid
export API_KEY="secret-123"

# Invalid - error if set to empty
export API_KEY=""
```

---

## Running the Example

```bash
cd examples/04_advanced_features
go run main.go
```

With custom environment:
```bash
export SERVER_URL=https://api.example.com
export BASE_PATH=/opt/myapp
export CONFIG_PATH='${BASE_PATH}/config'
export API_KEY=my-secret-key
export SERVICE_PORTS="database:5432,cache:6379"
export HEADERS="Accept=application/json"

go run main.go
```

---

## Key Takeaways

- **Variable expansion** helps compose configuration from other variables
- **URL parsing** provides type-safe URL handling
- **Maps** allow flexible key-value configuration
- **notEmpty** adds validation for required non-empty values
- **Custom separators** let you format environment variables as needed
