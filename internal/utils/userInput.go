package utils

import (
	"strings"
	"unicode"
)

// ParseIDList splits input by commas or whitespace and returns non-empty IDs.
func ParseIDList(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	var ids []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		ids = append(ids, part)
	}

	return ids
}
