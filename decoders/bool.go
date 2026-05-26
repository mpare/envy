package decoders

import (
	"fmt"
	"reflect"
	"strconv"
)

type boolDecoder struct{}

func (b *boolDecoder) CanDecode(field reflect.Value) bool {
	if implementsSelfDecoder(field) {
		return false
	}
	
	return field.Kind() == reflect.Bool
}

func (b *boolDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	parsedBool, err := strconv.ParseBool(raw)
	if err != nil {
		return fmt.Errorf("invalid bool %q: %w", raw, err)
	}

	field.SetBool(parsedBool)
	return
}
