package utils

import (
	"strings"
)

func Normalize(raw string, preset map[string]string, fallback string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	_, keyExists := preset[normalized]

	if keyExists {
		return normalized
	}
	return fallback
}

func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
