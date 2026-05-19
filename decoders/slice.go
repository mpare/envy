package decoders

import (
	"fmt"
	"reflect"
	"strings"
)

const defaultSeparator = ","

type sliceDecoder struct{}

func (s *sliceDecoder) CanDecode(field reflect.Value) bool {
	return field.Kind() == reflect.Slice
}

func (s *sliceDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	separator := tag.GetSeparator()
	if separator == "" {
		separator = defaultSeparator
	}

	parts := strings.Split(raw, separator)
	slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
	
	for i, part := range parts {
		if err := Decode(slice.Index(i), strings.TrimSpace(part), tag); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	field.Set(slice)
	return
}
