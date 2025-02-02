package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Custom validation function for mongo db object IDs
func IsValidMongoObjectID(id string) bool {
	// Ensure the ID is a 24-character hexadecimal string
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}

// convert str ids to mongo db object ids
func GetMongoObjectIds(ids []string) (objectIDs []primitive.ObjectID) {
	// Convert string IDs to ObjectIDs
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Error().Err(err).Msgf("Invalid ObjectID %s: %v", id, err)
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	return objectIDs
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

func ExcludeParamsFromStruct(data interface{}, params []string) (bson.M, error) {
	result := bson.M{}
	v := reflect.ValueOf(data)

	// Ensure we are working with a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// Loop through all fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		// Skip the field if match is found
		if slices.Contains(params, strings.ToLower(field.Name)) || field.Tag.Get("bson") == "_id" {
			continue
		}

		// Only include non-empty fields
		if !fieldValue.IsZero() {
			result[field.Tag.Get("bson")] = fieldValue.Interface()
		}
	}

	return result, nil
}
