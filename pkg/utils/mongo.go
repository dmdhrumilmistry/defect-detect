package utils

import "regexp"

// Custom validation function for mongo db object IDs
func isValidMongoObjectID(id string) bool {
	// Ensure the ID is a 24-character hexadecimal string
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}
