package decoders

import (
	"fmt"
	"reflect"
)

type TagReader interface {
	GetKey() string
	GetSeparator() string
	GetDefaultValue() string
	IsRequired() bool
	GetPrefix() string
}

// Decoder is for built-in/generic decoders that check multiple types.
type Decoder interface {
	CanDecode(field reflect.Value) bool
	Decode(field reflect.Value, raw string, tag TagReader) error
}

// SelfDecoder is for types that know how to decode themselves.
// Types implementing this interface can provide custom decoding logic.
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
