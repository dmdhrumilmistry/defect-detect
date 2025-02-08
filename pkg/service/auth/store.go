package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/dmdhrumilmistry/defect-detect/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
		log.Error().Err(err).Msg("failed to insert user")
		return "", err
	}

	return (result.InsertedID).(primitive.ObjectID).Hex(), nil
}

func (c *AuthStore) GetTotalCount(filter interface{}, collection *mongo.Collection) (int64, error) {
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

func (a *AuthStore) GetUserByEmail(email string, duration int) (user types.User, err error) {
	filter := bson.M{
		"email": email,
	}
	users, err := utils.GetObjectsUsingFilter[types.User](a.userCollection, filter, 1, 1, duration)
	if err != nil {
		return user, err
	}

	if len(users) == 0 {
		log.Warn().Msg("user not found for provided email id")
		return user, fmt.Errorf("user not found for provided email id")
	}

	if len(users) > 1 {
		log.Warn().Msgf("Multiple users found for provided email id")
	}
	user = users[0]

	return user, nil
}

// HasPermission checks if a user has access to a given resources (attributes).
func (c *AuthStore) HasPermission(user types.User, attributes []string, authOperator string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DefaultConfig.DbQueryTimeout)*time.Second)
	defer cancel()

	// // Fetch the user's groups from MongoDB
	// var user types.User
	// err := c.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Error fetching user")
	// 	return false, err
	// }

	// this is checked while validating JWT
	// check whether user is inactive
	// if !user.IsActive {
	// 	return false, fmt.Errorf("user is inactive")
	// }

	// return true is user is super user
	if user.IsSuperUser {
		return true, nil
	}

	// check if auth operator is valid else configure it to default AND
	var attributeFilterSymbol string
	switch strings.ToUpper(authOperator) {
	case "OR":
		attributeFilterSymbol = "$in"
	default: // AND
		attributeFilterSymbol = "$all"
	}

	if len(user.Groups) == 0 {
		log.Warn().Msgf("User %s does not belong to any group", user.Id)
		return false, nil
	}

	// Check if any of user's groups has the required attribute
	filter := bson.M{
		"_id":        bson.M{"$in": user.Groups},                // Find groups the user belongs to
		"attributes": bson.M{attributeFilterSymbol: attributes}, // Check if resource exists in attributes
	}

	log.Debug().Msgf("filter: %v", filter)

	count, err := c.groupCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching group attributes")
		return false, err
	}

	return count > 0, nil // If count > 0, user has permission
}

// middleware for validating JWT token
func (a *AuthStore) WithJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from request header
		tokenString := GetTokenFromRequest(c.Request)
		log.Info().Msgf("token: %s", tokenString)

		// validate token
		jwtToken, err := ValidateJWT(tokenString)
		if err != nil {
			log.Error().Err(err).Msg("failed to validate jwt token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
			return
		}

		if !jwtToken.Valid {
			log.Error().Err(err).Msg("invalid jwt")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
			return
		}

		// extract user id from token
		userId, ok := jwtToken.Claims.(jwt.MapClaims)[string(UserCtxKey)].(string)
		if !ok {
			log.Error().Err(err).Msg("failed to extract user Id from JWT token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
			return
		}

		// fetch user by id
		user, err := a.GetUserById(userId, config.DefaultConfig.DbQueryTimeout)
		if err != nil {
			log.Printf("error while fetching user: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "failed to fetch user details"})
			return
		}

		// check whether user is inactive
		if !user.IsActive {
			log.Warn().Msgf("inactive user tried to login: %s", user.Id)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "user is inactive"})
			return
		}

		// set user in context
		ctx := context.WithValue(c.Request.Context(), UserCtxKey, user)
		c.Request = c.Request.WithContext(ctx)

		// call handler function
		c.Next()
	}
}

// middleware for validating User permission
func (a *AuthStore) ValidatePerms(attributes []string, authOperator string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Request.Context().Value(UserCtxKey).(types.User)
		if !ok {
			log.Error().Msg("failed to retrieve user from context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "failed to check permissions"})
			return
		}

		hasAccess, err := a.HasPermission(user, attributes, authOperator)
		if err != nil {
			log.Error().Err(err).Msg("failed to check permissions")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "failed to check permissions"})
			return
		}

		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		// call handler function
		c.Next()
	}
}
