package decoders

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	intBase    = 10
	intBitSize = 64
)

type intDecoder struct{}

func (i *intDecoder) CanDecode(field reflect.Value) bool {
	return field.Kind() == reflect.Int ||
		field.Kind() == reflect.Int8 ||
		field.Kind() == reflect.Int16 ||
		field.Kind() == reflect.Int32 ||
		field.Kind() == reflect.Int64
}

func (i *intDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	parsedInt, err := strconv.ParseInt(raw, intBase, intBitSize)
	if err != nil {
		return fmt.Errorf("invalid int %q: %w", raw, err)
	}

	field.SetInt(parsedInt)
	return
}
