# Custom Decoder Example

This example demonstrates how to implement custom decoders for complex types that aren't natively supported by envy.

## Run

```bash
APP_NAME=api \
  FEATURES='{"auth":true,"logging":true,"caching":false,"ratelimit":true}' \
  METADATA='{"version":"1.0","env":"production","team":"platform"}' \
  go run main.go
```

## Output

```
Application: api

=== Features ===
  auth: true
  caching: false
  logging: true
  ratelimit: true

=== Metadata ===
  env: production
  team: platform
  version: 1.0
```

## Key Concepts

### SelfDecoder Interface

Any type can implement custom parsing by implementing the `SelfDecoder` interface:

```go
type SelfDecoder interface {
	Decode(field reflect.Value, raw string, tag decoders.TagReader) error
}
```

### Implementation Steps

1. **Define your type** - Create a custom type (e.g., `JSONData`)
2. **Implement Decode()** - Parse the raw string value and set the field
3. **Use in config** - Add fields with your custom type in configuration struct
4. **Automatic detection** - envy automatically uses your decoder when needed

### Common Use Cases

- **JSON objects** - Complex structured data
- **YAML configs** - Alternative configuration formats
- **Custom parsers** - Domain-specific parsing logic
- **Encrypted values** - Decrypt environment variable values
- **Remote configs** - Fetch configuration from external sources

## Error Handling

If JSON parsing fails:

```bash
APP_NAME=api \
  FEATURES='invalid-json' \
  go run main.go
```

envy will collect the error and report it with field and environment variable information.
