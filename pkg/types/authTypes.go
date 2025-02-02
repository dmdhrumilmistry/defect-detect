package types

import "time"

type AuthStore interface {
}

// Users can only be managed by groups as permissions to be only provided to groups
// instead of inidividual users
type User struct {
	Id        string    `json:"user_id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Groups    []string  `json:"group_ids" bson:"group_ids,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
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
