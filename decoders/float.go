package decoders

import (
	"fmt"
	"reflect"
	"strconv"
)

const floatBitSize = 64

type floatDecoder struct{}

func (f *floatDecoder) CanDecode(field reflect.Value) bool {
	return field.Kind() == reflect.Float32 || field.Kind() == reflect.Float64
}

func (f *floatDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	parsedFloat, err := strconv.ParseFloat(raw, floatBitSize)
	if err != nil {
		return fmt.Errorf("invalid float %q: %w", raw, err)
	}

	field.SetFloat(parsedFloat)
	return
}
