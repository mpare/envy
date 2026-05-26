package envy

import (
	"strings"
)

const (
	tagSeparator = ","
	requiredTag  = "required"
	defaultTag   = "default="
	separatorTag = "separator="
	prefixTag    = "prefix="
	expandTag    = "expand"
	fileTag      = "file"
	notEmptyTag  = "notEmpty"
	keyValSepTag = "keyValSeparator="
)

// Tag represents a parsed 'env' struct tag.
// It contains the environment variable name, default value, required flag, separator, prefix, and options.
type Tag struct {
	key             string
	defaultValue    string
	required        bool
	separator       string
	prefix          string
	expand          bool
	file            bool
	notEmpty        bool
	keyValSeparator string
}

func (t Tag) GetKey() string {
	return t.key
}

func (t Tag) GetSeparator() string {
	return t.separator
}

func (t Tag) GetDefaultValue() string {
	return t.defaultValue
}

func (t Tag) IsRequired() bool {
	return t.required
}

func (t Tag) GetPrefix() string {
	return t.prefix
}

func (t Tag) IsExpand() bool {
	return t.expand
}

func (t Tag) IsFile() bool {
	return t.file
}

func (t Tag) IsNotEmpty() bool {
	return t.notEmpty
}

func (t Tag) GetKeyValSeparator() string {
	return t.keyValSeparator
}

func parseTag(tag string) Tag {
	if tag == "" {
		return Tag{}
	}

	parts := strings.Split(tag, tagSeparator)
	result := Tag{key: strings.TrimSpace(parts[0])}

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)

		if part == requiredTag {
			result.required = true
			continue
		}

		if part == expandTag {
			result.expand = true
			continue
		}

		if part == fileTag {
			result.file = true
			continue
		}

		if part == notEmptyTag {
			result.notEmpty = true
			continue
		}

		if strings.HasPrefix(part, defaultTag) {
			result.defaultValue = strings.TrimPrefix(part, defaultTag)
			continue
		}

		if strings.HasPrefix(part, separatorTag) {
			result.separator = strings.TrimPrefix(part, separatorTag)
			continue
		}

		if strings.HasPrefix(part, keyValSepTag) {
			result.keyValSeparator = strings.TrimPrefix(part, keyValSepTag)
			continue
		}

		if strings.HasPrefix(part, prefixTag) {
			result.prefix = strings.TrimPrefix(part, prefixTag)
			continue
		}
	}

	return result
}
