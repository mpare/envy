# Basic Configuration Example

This example demonstrates the simplest way to use envy - loading basic configuration types.

## Run

```bash
# Using defaults
go run main.go

# Override with environment variables
PORT=3000 DEBUG=true go run main.go

# Or set them all
APP_NAME=myservice PORT=9000 DEBUG=true go run main.go
```

## Configuration

- `APP_NAME` - Application name (default: `myapp`)
- `PORT` - Server port (default: `8080`)
- `DEBUG` - Enable debug mode (default: `false`)

## Key Concepts

- **Default values**: Use `default=value` tag to provide fallback values
- **Type conversion**: Integers and booleans are automatically converted
- **Simple fields**: Each field maps to a single environment variable
