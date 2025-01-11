package types

type ProjectHandler interface {
}

type Project struct {
	Id            string   `json:"id" bson:"_id,omitempty"`
	Name          string   `json:"name" bson:"name" binding:"required"`
	Description   string   `json:"description" bson:"description" binding:"required"`
	SbomsToRetain int      `json:"sboms_to_retain" bson:"sboms_to_retain" binding:"required"`
	Sboms         []string `json:"sboms" bson:"sboms"`
	Links         []string `json:"links" bson:"links" binding:"min=0,max=5"`
}
