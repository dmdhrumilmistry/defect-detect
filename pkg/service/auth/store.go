package auth

import (
	"context"
	"strings"
	"time"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USER_COLLECTION = "user"
const GROUP_COLLECTION = "group"

type AuthStore struct {
	db              *mongo.Database
	userCollection  *mongo.Collection
	groupCollection *mongo.Collection
}

func NewAuthStore(db *mongo.Database) *AuthStore {
	return &AuthStore{
		db:              db,
		userCollection:  db.Collection(USER_COLLECTION),
		groupCollection: db.Collection(GROUP_COLLECTION),
	}
}

func (a *AuthStore) CreateUser(user types.User) (string, error) {
	result, err := a.userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert sbom")
		return "", err
	}

	return (result.InsertedID).(primitive.ObjectID).Hex(), nil
}

// Collection Type: user/group. Default user collection
func (c *AuthStore) GetTotalCount(filter interface{}, collectionType string) (int64, error) {
	// configure collection as per collection type defaults to user collection
	collection := c.userCollection
	if collectionType == "group" {
		collection = c.groupCollection
	}

	// Get total count of documents
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (c *AuthStore) GetUserById(idParam string, duration int) (types.User, error) {
	var object types.User

	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return object, err
	}

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	err = c.userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&object)
	if err != nil {
		return object, err
	}

	return object, nil
}

// HasPermission checks if a user has access to a given resources (attributes).
func (c *AuthStore) HasPermission(userID, attributes []string, authOperator string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DefaultConfig.DbQueryTimeout)*time.Second)
	defer cancel()

	// Fetch the user's groups from MongoDB
	var user types.User
	err := c.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching user")
		return false, err
	}

	// check if auth operator is valid else configure it to default AND
	var attributeFilterSymbol string
	switch strings.ToUpper(authOperator) {
	case "OR":
		attributeFilterSymbol = "$in"
	default: // AND
		attributeFilterSymbol = "$all"
	}

	// Check if any of user's groups has the required attribute
	filter := bson.M{
		"_id":        bson.M{"$in": user.Groups},                // Find groups the user belongs to
		"attributes": bson.M{attributeFilterSymbol: attributes}, // Check if resource exists in attributes
	}

	count, err := c.groupCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching group attributes")
		return false, err
	}

	return count > 0, nil // If count > 0, user has permission
}
