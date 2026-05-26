# Validation and Error Handling Example

This example demonstrates envy's validation features and how to handle validation errors.

## Validation Features

### 1. Required Fields (`required` tag)

Mark a field as required - it must be set in the environment or envy returns an error.

```go
AppName string `env:"APP_NAME,required"`
```

**Behavior:**
- If `APP_NAME` is not set → validation error
- If `APP_NAME` is set to empty string `""` → OK (field is set)
- If `APP_NAME` is set to any value → OK

---

### 2. Not Empty Validation (`notEmpty` tag)

Validate that a field is not an empty string.

```go
APIKey string `env:"API_KEY,notEmpty"`
```

**Behavior:**
- If `API_KEY` is not set → OK (uses zero value)
- If `API_KEY=""` (empty) → validation error
- If `API_KEY` has any value → OK

**Difference from `required`:**
- `required`: Must be set (any value, including empty string)
- `notEmpty`: Must not be empty (if set, must have content)

---

### 3. Default Values (`default=` tag)

Provide a fallback value if the environment variable is not set.

```go
Port int `env:"PORT,default=8080"`
```

**Behavior:**
- If `PORT` not set → uses `8080`
- If `PORT` set to `9000` → uses `9000`

---

## Error Handling

### Validation Error Type

All validation errors are collected into a `ValidationError` that contains multiple `FieldError` items:

```go
type ValidationError struct {
    Errors []FieldError
}

type FieldError struct {
    Field   string  // Struct field name
    EnvKey  string  // Environment variable name
    Message string  // Error message
}
```

### Checking Errors

```go
if err := envy.Load(&cfg); err != nil {
    var ve *envy.ValidationError
    if errors.As(err, &ve) {
        for _, fieldErr := range ve.Errors {
            fmt.Printf("%s (%s): %s\n", fieldErr.Field, fieldErr.EnvKey, fieldErr.Message)
        }
    }
}
```

### Error Examples

**Missing required field:**
```
AppName (APP_NAME): required but not set
```

**Empty value with notEmpty:**
```
APIKey (API_KEY): must not be empty
```

**Invalid type conversion:**
```
Port (PORT): invalid int "abc": strconv.ParseInt: parsing "abc": invalid syntax
```

**Invalid duration:**
```
Timeout (TIMEOUT): invalid duration "invalid": time.ParseDuration: ...
```

---

## Example Scenarios

### Scenario 1: Valid Configuration

All required fields set with valid values → Success

```go
validEnv := map[string]string{
    "APP_NAME":      "myapp",
    "API_KEY":       "secret-key",
    "DATABASE_URL":  "postgres://localhost/mydb",
}

err := envy.LoadFrom(&cfg, validEnv)
// err == nil ✓
```

### Scenario 2: Missing Required Field

```go
invalidEnv := map[string]string{
    "APP_NAME": "myapp",
    // DATABASE_URL missing!
}

err := envy.LoadFrom(&cfg, invalidEnv)
// err != nil ✗
// Message: DatabaseURL (DATABASE_URL): required but not set
```

### Scenario 3: Empty Value with notEmpty

```go
invalidEnv := map[string]string{
    "APP_NAME": "myapp",
    "API_KEY":  "",  // Empty! Violates notEmpty
}

err := envy.LoadFrom(&cfg, invalidEnv)
// err != nil ✗
// Message: APIKey (API_KEY): must not be empty
```

### Scenario 4: Multiple Validation Errors

All errors are collected and reported together:

```go
invalidEnv := map[string]string{
    // APP_NAME missing
    "API_KEY": "",
    // DATABASE_URL missing
}

err := envy.LoadFrom(&cfg, invalidEnv)
// err != nil with 3 field errors:
// 1. AppName (APP_NAME): required but not set
// 2. APIKey (API_KEY): must not be empty
// 3. DatabaseURL (DATABASE_URL): required but not set
```

**Benefits:**
- See all problems at once instead of fixing one and retrying
- Better user experience for configuration errors
- Easier debugging of configuration issues

---

## Running the Example

```bash
cd examples/06_validation_and_errors
go run main.go
```

The example will demonstrate:
1. **Valid config** - successfully loaded
2. **Missing required field** - validation error
3. **Empty value** - notEmpty validation error
4. **Multiple errors** - all collected and displayed

---

## Best Practices

### 1. **Use `required` for Critical Fields**

```go
type DatabaseConfig struct {
    Host     string `env:"HOST,required"`
    User     string `env:"USER,required"`
    Password string `env:"PASSWORD,required"`
}
```

### 2. **Use `notEmpty` When Presence Isn't Enough**

```go
type APIConfig struct {
    // The key must not just exist, but must have a value
    APIKey string `env:"API_KEY,notEmpty"`
}
```

### 3. **Provide Sensible Defaults**

```go
type ServerConfig struct {
    Port    int    `env:"PORT,default=8080"`
    LogLevel string `env:"LOG_LEVEL,default=info"`
}
```

### 4. **Handle Errors Appropriately**

```go
// In main() or init:
var cfg Config
if err := envy.Load(&cfg); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}

// In tests:
var cfg Config
if err := envy.LoadFrom(&cfg, testEnv); err != nil {
    t.Fatalf("Failed to load test config: %v", err)
}
```

### 5. **Document Validation Rules**

```go
type Config struct {
    // AppName is required and must not be empty
    AppName string `env:"APP_NAME,required,notEmpty"`

    // APIKey is required and must not be empty
    APIKey string `env:"API_KEY,notEmpty"`

    // Port is optional, defaults to 8080 if not set
    Port int `env:"PORT,default=8080"`
}
```

---

## Validation Rule Combinations

| Tag | Required | Empty OK | If Missing | If `=""` |
|-----|----------|----------|-----------|---------|
| `env:"X"` | ❌ | ✓ | Zero value | Empty |
| `env:"X,required"` | ✓ | ✓ | Error | Empty |
| `env:"X,notEmpty"` | ❌ | ❌ | Zero value | Error |
| `env:"X,required,notEmpty"` | ✓ | ❌ | Error | Error |
| `env:"X,default=val"` | ❌ | ✓ | `val` | Empty |
| `env:"X,required,default=val"` | ✓ | ✓ | `val` | `val` |

---

## Key Takeaways

- Use `required` for mandatory fields
- Use `notEmpty` when the field must have a value
- Use `default` to provide fallback values
- All validation errors are collected together
- Check `ValidationError.Errors` for detailed field errors
- Proper error handling improves user experience
