# Examples Summary

## Overview

The examples directory contains 6 comprehensive examples demonstrating all features of the envy library, from basic usage to advanced features.

## Example Breakdown

| # | Name | Focus | When to Use |
|---|------|-------|------------|
| 01 | [Basic](01_basic) | Simple types, defaults | Getting started |
| 02 | [Microservice](02_microservice) | Nested structs, prefixes | Multi-config services |
| 03 | [Custom Decoder](03_custom_decoder) | SelfDecoder interface, JSON | Complex types |
| 04 | [Advanced Features](04_advanced_features) | Expansion, URLs, maps | New features |
| 05 | [File-Based Config](05_file_based_config) | File reading | Security, large configs |
| 06 | [Validation & Errors](06_validation_and_errors) | Validation, error handling | Error management |

---

## Quick Navigation

### I'm New to Envy
Start with [01_basic](01_basic) to understand the fundamentals.

### I'm Building a Microservice
See [02_microservice](02_microservice) for realistic configuration patterns.

### I Need to Parse Complex Types
Check [03_custom_decoder](03_custom_decoder) for the SelfDecoder pattern.

### I Want to Use New Features
Explore [04_advanced_features](04_advanced_features) for URL parsing, maps, and variable expansion.

### I Need Secure Configuration
Look at [05_file_based_config](05_file_based_config) for file-based secrets.

### I Need to Handle Errors
Study [06_validation_and_errors](06_validation_and_errors) for comprehensive error handling.

---

## Running Examples

Each example can be run independently:

```bash
cd examples/01_basic
go run main.go
```

With environment variables:

```bash
cd examples/02_microservice
APP_NAME=myservice DB_HOST=localhost DB_USER=admin DB_PASSWORD=secret go run main.go
```

---

## Feature Comparison Matrix

| Feature | 01 | 02 | 03 | 04 | 05 | 06 |
|---------|----|----|----|----|----|----|
| Basic types | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Default values | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Required fields | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ |
| Nested structs | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Custom decoders | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Variable expansion | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| URL parsing | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| Map support | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| File reading | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| Empty validation | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Error handling | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## Key Takeaways

### 1. Type Safety
Envy ensures environment variables are properly typed - no more string parsing in your code.

### 2. Error Aggregation
All validation errors are collected in one pass, giving you complete visibility into configuration issues.

### 3. Flexibility
From simple types to custom JSON/YAML, envy handles various configuration patterns.

### 4. Security
File reading and environment variable patterns support secure handling of secrets.

### 5. Real-World Patterns
Examples follow production-grade patterns used in actual microservices.

---

## Common Use Cases

### API Server Configuration
```
See 02_microservice + 04_advanced_features
```

### Database Connection Strings
```
See 02_microservice + 05_file_based_config
```

### Feature Flags and Metadata
```
See 03_custom_decoder
```

### Multi-Environment Deployment
```
Combine examples with different env files
```

### Docker/Kubernetes Secrets
```
See 05_file_based_config for secret mounting patterns
```

---

## Pro Tips

1. **Use Prefixes for Organization** (Example 02)
   - Group related configs with prefixes (DB_, CACHE_, etc.)

2. **Leverage Type Conversion** (All Examples)
   - Let envy handle int, bool, duration conversions

3. **Combine Tags** (Example 06)
   - Use `required,notEmpty,default` together for validation

4. **Test with LoadFrom()** (All Examples)
   - Pass custom env maps for unit tests

5. **Use File Reading for Secrets** (Example 05)
   - Keeps sensitive data out of environment variables

---

## Learning Path

```
01_basic → 02_microservice → 03_custom_decoder → 
04_advanced_features → 05_file_based_config → 06_validation_and_errors
```

Each example builds on concepts from previous ones while introducing new capabilities.

---

## Documentation Links

- [Main README](../README.md) - Library overview
- [API Documentation](https://pkg.go.dev/github.com/mpare/envy) - Full API reference
- Each example's README - Detailed explanation of features
