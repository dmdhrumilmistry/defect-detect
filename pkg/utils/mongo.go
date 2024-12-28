package utils

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

// Custom validation function for mongo db object IDs
func isValidMongoObjectID(id string) bool {
	// Ensure the ID is a 24-character hexadecimal string
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}

func BuildDynamicContainsFilter(conditions map[string][]string) bson.M {
	filter := bson.M{}

	for field, values := range conditions {
		if len(values) > 0 {
			filter[field] = bson.M{
				"$in": values,
			}
		}
	}

	return filter
}

func RemoveEmptyStrings(slice []string) []string {
	var result []string
	for _, str := range slice {
		if str != "" { // Only add non-empty strings
			result = append(result, str)
		}
	}
	return result
}
