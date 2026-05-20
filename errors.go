package envy

import (
	"fmt"
	"strings"
)

// ValidationError represents one or more field validation errors during environment variable loading.
// It aggregates all field errors so users can see all problems at once instead of failing on the first error.
type ValidationError struct {
	Errors []FieldError
}

// FieldError describes a single field validation error.
type FieldError struct {
	Field   string
	EnvKey  string
	Message string
}

func (v *ValidationError) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d env error(s):\n", len(v.Errors))

	for _, error := range v.Errors {
		fmt.Fprintf(&builder, "  - %s (%s): %s\n", error.Field, error.EnvKey, error.Message)
	}

	return builder.String()
}
