// Package decoders provides type conversion logic for environment variable values.
// It includes built-in decoders for common types (string, int, float, bool, duration, slices)
// and support for custom decoders through the SelfDecoder interface.
package decoders

import (
	"fmt"
	"reflect"
)

// TagReader provides access to parsed environment variable tag options.
type TagReader interface {
	GetKey() string
	GetSeparator() string
	GetDefaultValue() string
	IsRequired() bool
	GetPrefix() string
}

// Decoder defines the interface for built-in type decoders.
// Implementations check if they can decode a specific field type, then decode raw string values.
type Decoder interface {
	CanDecode(field reflect.Value) bool
	Decode(field reflect.Value, raw string, tag TagReader) error
}

// SelfDecoder defines the interface for types that can decode themselves from environment variable strings.
// This allows custom types (e.g., JSON objects, YAML configs) to provide their own parsing logic.
//
// Example implementation:
//
//	type JSONData map[string]interface{}
//
//	func (j *JSONData) Decode(field reflect.Value, raw string, tag TagReader) error {
//		return json.Unmarshal([]byte(raw), j)
//	}
type SelfDecoder interface {
	Decode(field reflect.Value, raw string, tag TagReader) error
}

var decoders = []Decoder{
	&durationDecoder{},
	&stringDecoder{},
	&intDecoder{},
	&floatDecoder{},
	&boolDecoder{},
	&sliceDecoder{},
}

func Decode(field reflect.Value, raw string, tag TagReader) error {
	for _, decoder := range decoders {
		if decoder.CanDecode(field) {
			return decoder.Decode(field, raw, tag)
		}
	}

	if decoderValue := field.Addr(); decoderValue.Type().Implements(reflect.TypeOf((*SelfDecoder)(nil)).Elem()) {
		if decoder, ok := decoderValue.Interface().(SelfDecoder); ok {
			return decoder.Decode(field, raw, tag)
		}
	}

	return fmt.Errorf("unsupported type: %s (no decoder found)", field.Type())
}
