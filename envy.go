package envy

import (
	"fmt"
	"os"
	"reflect"
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

func MustLoad(destination any) {
	if err := Load(destination); err != nil {
		panic(err)
	}
}

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
