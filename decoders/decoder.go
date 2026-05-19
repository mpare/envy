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

type Decoder interface {
	CanDecode(field reflect.Value) bool
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

	return fmt.Errorf("unsupported type: %s", field.Kind())
}
