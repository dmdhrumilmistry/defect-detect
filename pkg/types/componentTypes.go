package types

type ComponentStore interface {
	AddComponentUsingSbom(sbom Sbom) ([]string, error)
	GetComponentTotalCount() (int64, error)
	GetPaginatedComponents(page, limit, duration int) ([]Component, error)
	GetComponentById(idParam string, duration int) ([]Component, error)
	GetComponentByName(name string, duration int) ([]Component, error)
}

type Component struct {
	Id               string   `json:"component_id" bson:"_id,omitempty"`
	Name             string   `json:"name" bson:"name"`
	Version          string   `json:"version" bson:"version"`
	PackageUrl       string   `json:"purl" bson:"purl"`
	Licenses         []string `json:"licenses" bson:"licenses"`
	Type             string   `json:"type" bson:"type"`
	ComponentName    string   `json:"component_name" bson:"component_name"`
	ComponentVersion string   `json:"component_version" bson:"component_version"`
	// TODO: add vuln field in component
}
