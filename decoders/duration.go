package decoders

import (
	"fmt"
	"reflect"
	"time"
)

type durationDecoder struct{}

func (d *durationDecoder) CanDecode(field reflect.Value) bool {
	return field.Type() == reflect.TypeOf(time.Duration(0))
}

func (d *durationDecoder) Decode(field reflect.Value, raw string, tag TagReader) (err error) {
	parsedDuration, err := time.ParseDuration(raw)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", raw, err)
	}

	field.SetInt(int64(parsedDuration))
	return
}
