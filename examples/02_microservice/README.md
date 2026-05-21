# Microservice Configuration Example

This example demonstrates how to use envy for a realistic microservice setup with nested configuration structs and environment variable prefixes.

## Run

```bash
# Using minimum required variables
DB_HOST=localhost DB_USER=admin DB_PASSWORD=secret APP_NAME=my-service go run main.go

# Full example
APP_NAME=api-service \
  PORT=3000 \
  DEBUG=true \
  LOG_LEVEL=debug \
  DB_HOST=db.example.com \
  DB_PORT=5432 \
  DB_USER=dbuser \
  DB_PASSWORD=dbpass \
  DB_NAME=mydb \
  CACHE_HOST=cache.example.com \
  CACHE_PORT=6379 \
  CACHE_TTL=2h \
  go run main.go
```

## Configuration

### Application Settings
- `APP_NAME` - Service name (**required**)
- `PORT` - Server port (default: `8080`)
- `DEBUG` - Enable debug mode (default: `false`)
- `LOG_LEVEL` - Logging level (default: `info`)
- `TIMEOUT` - Request timeout (default: `30s`)

### Database Settings (prefix: `DB_`)
- `DB_HOST` - Database host (**required**)
- `DB_PORT` - Database port (default: `5432`)
- `DB_USER` - Database user (**required**)
- `DB_PASSWORD` - Database password (**required**)
- `DB_NAME` - Database name (default: `postgres`)

### Cache Settings (prefix: `CACHE_`)
- `CACHE_HOST` - Cache host (default: `localhost`)
- `CACHE_PORT` - Cache port (default: `6379`)
- `CACHE_TTL` - Cache TTL (default: `1h`)
- `CACHE_ENABLED` - Enable cache (default: `true`)

## Key Concepts

- **Nested structs**: Group related configuration in nested structs
- **Prefixes**: Use `prefix=PREFIX_` to organize environment variables by component
- **Required fields**: Use `required` tag to enforce mandatory configuration
- **Type safety**: Durations, integers, and booleans are automatically converted
- **Defaults**: Non-required fields have sensible defaults
