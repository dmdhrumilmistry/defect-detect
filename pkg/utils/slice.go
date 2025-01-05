package utils

import (
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

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

// finds and returns first element matching provided pattern
// returns empty str if no match found or any err occurs
func FindRegexMatchEle(pattern string, s []string) string {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Error().Err(err).Msg("failed to compile regex")
		return ""
	}
	for _, ele := range s {
		if regex.MatchString(ele) {
			return ele
		}
	}

	return ""
}
