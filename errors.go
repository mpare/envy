package envy

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors []FieldError
}

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
