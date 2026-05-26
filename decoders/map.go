package decoders

import (
	"fmt"
	"reflect"
	"strings"
)

type mapDecoder struct{}

func (m *mapDecoder) CanDecode(field reflect.Value) bool {
	if implementsSelfDecoder(field) {
		return false
	}
	
	return field.Kind() == reflect.Map
}

func (m *mapDecoder) Decode(field reflect.Value, raw string, tag TagReader) error {
	separator := tag.GetSeparator()
	if separator == "" {
		separator = defaultSeparator
	}

	keyValSeparator := tag.GetKeyValSeparator()
	if keyValSeparator == "" {
		keyValSeparator = ":"
	}

	mapType := field.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)

	if raw == "" {
		field.Set(newMap)
		return nil
	}

	entries := strings.Split(raw, separator)
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, keyValSeparator, 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid map entry format: %q (expected key%svalue)", entry, keyValSeparator)
		}

		keyStr := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		keyVal := reflect.New(keyType).Elem()
		if err := Decode(keyVal, keyStr, tag); err != nil {
			return fmt.Errorf("invalid map key: %w", err)
		}

		valueVal := reflect.New(valueType).Elem()
		if err := Decode(valueVal, valueStr, tag); err != nil {
			return fmt.Errorf("invalid map value for key %q: %w", keyStr, err)
		}

		newMap.SetMapIndex(keyVal, valueVal)
	}

	field.Set(newMap)
	return nil
}
