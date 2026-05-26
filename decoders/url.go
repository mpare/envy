package decoders

import (
	"fmt"
	"net/url"
	"reflect"
)

type urlDecoder struct{}

func (u *urlDecoder) CanDecode(field reflect.Value) bool {
	if implementsSelfDecoder(field) {
		return false
	}
	
	return field.Type() == reflect.TypeOf((*url.URL)(nil)).Elem()
}

func (u *urlDecoder) Decode(field reflect.Value, raw string, tag TagReader) error {
	parsedUrl, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	field.Set(reflect.ValueOf(*parsedUrl))
	return nil
}
