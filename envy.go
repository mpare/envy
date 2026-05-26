// Package envy provides type-safe environment variable loading with validation, defaults, and nested struct support.
// Zero external dependencies, uses Go stdlib only.
//
// Envy allows you to define typed configuration structs with env tags and automatically load
// environment variables into them with full type conversion, validation, and error reporting.
//
// Basic usage:
//
//	type Config struct {
//		Port     int    `env:"PORT,default=8080"`
//		Debug    bool   `env:"DEBUG,default=false"`
//		Database string `env:"DATABASE_URL,required"`
//	}
//
//	var cfg Config
//	envy.MustLoad(&cfg)
//
// For more information, see https://github.com/mpare/envy
package envy

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/mpare/envy/decoders"
)

const (
	envTagKey         = "env"
	keyValueSeparator = "="
	firstChar         = 'A'
	lastChar          = 'Z'
	underscore        = '_'
)

var expandVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// Load reads environment variables from os.Environ() and populates the given struct.
// The struct must have exported fields with 'env' tags specifying environment variable names.
// Returns ValidationError if any required fields are missing or values cannot be parsed.
func Load(destination any) error {
	environ := os.Environ()
	envMap := make(map[string]string, len(environ))

	for _, env := range environ {
		parts := strings.SplitN(env, keyValueSeparator, 2)

		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	return LoadFrom(destination, envMap)
}

// MustLoad is like Load but panics if an error occurs.
// Useful for application startup where configuration errors should be fatal.
func MustLoad(destination any) {
	if err := Load(destination); err != nil {
		panic(err)
	}
}

// LoadFrom reads environment variables from the given map and populates the struct.
// This is useful for testing, allowing you to inject custom environment variable maps.
//
// Supported field types:
//   - string, int8, int16, int32, int64, int
//   - float32, float64
//   - bool (accepts: true, false, 1, 0, t, f)
//   - time.Duration
//   - url.URL
//   - slices: []string, []int, []float64, []bool (with configurable separator)
//   - maps: map[K]V where K and V are supported types (with configurable separators)
//   - nested structs (with configurable prefix)
//   - custom types implementing SelfDecoder interface
//
// Struct field tags:
//   - env:"NAME" - required; sets environment variable name
//   - env:"NAME,default=value" - optional with default value
//   - env:"NAME,required" - marks field as required
//   - env:"NAME,notEmpty" - field must not be empty if set
//   - env:"NAME,expand" - expand ${VAR} references in the value
//   - env:"NAME,file" - treat value as file path and read file content
//   - env:"NAME,separator=;" - for slices/maps, sets item separator (default: ",")
//   - env:"NAME,keyValSeparator=:" - for maps, sets key:value separator (default: ":")
//   - env:",prefix=PREFIX_" - for nested structs, sets environment variable prefix
func LoadFrom(destination any, envMap map[string]string) error {
	destinationValue := reflect.ValueOf(destination)
	if destinationValue.Kind() != reflect.Ptr || destinationValue.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer to a struct")
	}

	destinationValue = destinationValue.Elem()
	if destinationValue.Kind() != reflect.Struct {
		return fmt.Errorf("destination must point to a struct, got %s", destinationValue.Kind())
	}

	destinationType := destinationValue.Type()
	var errs []FieldError

	for i := 0; i < destinationType.NumField(); i++ {
		field := destinationType.Field(i)
		fieldValue := destinationValue.Field(i)

		if !field.IsExported() {
			continue
		}

		tagStr := field.Tag.Get(envTagKey)

		if tagStr == "" {
			if field.Type.Kind() != reflect.Struct {
				continue
			}
		}

		tag := parseTag(tagStr)

		if field.Type.Kind() == reflect.Struct {
			prefix := tag.prefix

			if prefix == "" && tagStr != "" {
				prefix = tag.key
			}

			if err := loadNestedStruct(fieldValue, envMap, prefix); err != nil {
				if ve, ok := err.(*ValidationError); ok {
					errs = append(errs, ve.Errors...)
				} else {
					errs = append(errs, FieldError{
						Field:   field.Name,
						EnvKey:  tag.key,
						Message: err.Error(),
					})
				}
			}

			continue
		}

		envKey := tag.key
		rawValue, exists := envMap[envKey]

		if !exists {
			if tag.required {
				errs = append(errs, FieldError{
					Field:   field.Name,
					EnvKey:  envKey,
					Message: "required but not set",
				})
			} else if tag.defaultValue != "" {
				rawValue = tag.defaultValue
			} else {
				continue
			}
		}

		if tag.expand {
			rawValue = expandVars(rawValue, envMap)
		}

		if tag.notEmpty && rawValue == "" {
			errs = append(errs, FieldError{
				Field:   field.Name,
				EnvKey:  envKey,
				Message: "must not be empty",
			})

			continue
		}

		if tag.file {
			content, err := readFileContent(rawValue)
			if err != nil {
				errs = append(errs, FieldError{
					Field:   field.Name,
					EnvKey:  envKey,
					Message: err.Error(),
				})

				continue
			}

			rawValue = content
		}

		if err := decoders.Decode(fieldValue, rawValue, tag); err != nil {
			errs = append(errs, FieldError{
				Field:   field.Name,
				EnvKey:  envKey,
				Message: err.Error(),
			})
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}

	return nil
}

func loadNestedStruct(structVal reflect.Value, envMap map[string]string, prefix string) error {
	if structVal.Kind() != reflect.Struct {
		return nil
	}

	structType := structVal.Type()
	var errs []FieldError

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structVal.Field(i)

		if !field.IsExported() {
			continue
		}

		tagStr := field.Tag.Get(envTagKey)
		tag := parseTag(tagStr)
		envKey := prefix + tag.key

		if tag.key == "" {
			envKey = prefix + toEnvName(field.Name)
		}

		rawValue, exists := envMap[envKey]

		if !exists {
			if tag.required {
				errs = append(errs, FieldError{
					Field:   field.Name,
					EnvKey:  envKey,
					Message: "required but not set",
				})
			} else if tag.defaultValue != "" {
				rawValue = tag.defaultValue
			} else {
				continue
			}
		}

		if tag.expand {
			rawValue = expandVars(rawValue, envMap)
		}

		if tag.notEmpty && rawValue == "" {
			errs = append(errs, FieldError{
				Field:   field.Name,
				EnvKey:  envKey,
				Message: "must not be empty",
			})

			continue
		}

		if tag.file {
			content, err := readFileContent(rawValue)
			if err != nil {
				errs = append(errs, FieldError{
					Field:   field.Name,
					EnvKey:  envKey,
					Message: err.Error(),
				})
				continue
			}

			rawValue = content
		}

		if err := decoders.Decode(fieldValue, rawValue, tag); err != nil {
			errs = append(errs, FieldError{
				Field:   field.Name,
				EnvKey:  envKey,
				Message: err.Error(),
			})
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}

	return nil
}

func expandVars(value string, envMap map[string]string) string {
	return expandVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		varName := match[2 : len(match)-1]
		if val, exists := envMap[varName]; exists {
			return val
		}

		return match
	})
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", path, err)
	}
	
	return string(content), nil
}

func toEnvName(value string) string {
	var result strings.Builder

	for i, char := range value {
		if i > 0 && char >= firstChar && char <= lastChar {
			result.WriteRune(underscore)
		}

		result.WriteRune(char)
	}

	return strings.ToUpper(result.String())
}
