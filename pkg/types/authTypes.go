package types

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthStore interface {
	CreateUser(user User) (string, error)
	GetTotalCount(filter interface{}, collection *mongo.Collection) (int64, error)
	GetUserById(idParam string, duration int) (User, error)
	GetUserByEmail(email string, duration int) (user User, err error)

	// Permissions
	WithJwtAuth() gin.HandlerFunc
	HasPermission(user User, attributes []string, authOperator string) (bool, error)
}

// Users can only be managed by groups as permissions to be only provided to groups
// instead of inidividual users
type User struct {
	Id          string    `json:"user_id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	AvatarUrl   string    `json:"avatar_url" bson:"avatar_url"`
	Groups      []string  `json:"group_ids" bson:"group_ids,omitempty"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	IsSuperUser bool      `json:"is_superuser" bson:"is_superuser"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

// Permissions should be only provided to groups
type Group struct {
	Id         string    `json:"group_id" bson:"_id,omitempty"`
	Name       string    `json:"name" bson:"name"`
	Users      []string  `json:"user_ids" bson:"user_ids,omitempty"`
	Attributes []string  `json:"attributes" bson:"attributes,omitempty"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}
