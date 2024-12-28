package utils

import "strings"

// removes empty string from slice after separating
func Split(s string, sep string) []string {
	parts := strings.Split(s, sep)
	writeIdx := 0

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			parts[writeIdx] = part
			writeIdx++
		}
	}

	// Slice down to only the valid elements
	return parts[:writeIdx]
}
