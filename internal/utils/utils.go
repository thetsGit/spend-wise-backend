package utils

import (
	"strings"
)

func Normalize(raw string, preset map[string]bool, fallback string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if preset[normalized] {
		return normalized
	}
	return fallback
}

func Keys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
