package types

type ProjectStore interface {
	AddProject(project Project) (string, error)
	GetTotalCount(filter interface{}) (int64, error)
	GetUsingFilter(filter interface{}, page, limit, duration int) ([]Project, error)
	GetProjectById(idParam string, duration int) ([]Project, error)
	GetByName(name string, duration int) ([]Project, error)
	UpdateById(payload Project, duration int) error
	DeleteByIds(idParams []string, duration int) (int64, error)
	DeleteById(idParam string, duration int) (int64, error)
	ValidateIds(ids []string) error
}

type Project struct {
	Id            string   `json:"id" bson:"_id,omitempty"`
	Name          string   `json:"name" bson:"name" binding:"required"`
	Description   string   `json:"description" bson:"description" binding:"required"`
	SbomsToRetain int      `json:"sboms_to_retain" bson:"sboms_to_retain" binding:"required,min=1,max=10"`
	Sboms         []string `json:"sboms" bson:"sboms"`
	Links         []string `json:"links" bson:"links" binding:"min=0,max=5"`
}
