package utils

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetObjectsUsingFilter[T any](collection *mongo.Collection, filter interface{}, page, limit, duration int) ([]T, error) {
	var objects []T

	// Calculate skip
	skip := (page - 1) * limit

	// MongoDB query options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return objects, err
	}
	defer cursor.Close(ctx)

	// Parse results
	if err := cursor.All(ctx, &objects); err != nil {
		return objects, err
	}

	return objects, nil
}
