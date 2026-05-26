package decoders

import (
	"reflect"
)

type stringDecoder struct{}

func (s *stringDecoder) CanDecode(field reflect.Value) bool {
	if implementsSelfDecoder(field) {
		return false
	}
	
	return field.Kind() == reflect.String
}

func (s *stringDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	field.SetString(raw)
	return
}
