package types

import (
	"time"

	"github.com/dmdhrumilmistry/m-paf/pkg/socketdev"
)

type ComponentStore interface {
	AddComponentUsingSbom(sbom Sbom) ([]string, error)
	GetComponentTotalCount(filter interface{}) (int64, error)
	GetPaginatedComponents(page, limit, duration int) ([]Component, error)
	GetComponentById(idParam string, duration int) ([]Component, error)
	GetComponentByName(name string, duration int) ([]Component, error)
	GetVulnerableComponents(componentNames, componentVersions, sbomIds, compTypes, compNames, purls, versions []string, page, limit, duration int) (components []Component, total int64, err error)
	DeleteByIds(idParams []string, param string, duration int) (int64, error)
	DeleteById(idParam string, param string, duration int) (int64, error)
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
	SbomId           string   `json:"sbom_id" bson:"sbom_id"`

	// Analyzers
	Vulns []Vuln `json:"vulns" bson:"vulns"`

	// M-Paf Analyzer
	PackageInfos []PackageInfo `json:"package_infos,omitempty"`
	// Alerts       []socketdev.Alert      `json:"alerts,omitempty"`
	// Scores       socketdev.Scores       `json:"scores,omitempty"`
	// Capabilities socketdev.Capabilities `json:"capabilities,omitempty"`
}

type PackageInfo struct {
	// ID        string `json:"id"`
	// Type      string `json:"type"`
	// Name      string `json:"name"`
	// Namespace string `json:"namespace"`
	// Files          string                 `json:"files"`
	// Version        string                 `json:"version"`
	// Qualifiers     socketdev.Qualifiers   `json:"qualifiers"`
	Scores         socketdev.Scores       `json:"scores"`
	Capabilities   socketdev.Capabilities `json:"capabilities"`
	License        string                 `json:"license"`
	Size           int                    `json:"size"`
	State          string                 `json:"state"`
	Alerts         []socketdev.AlertType  `json:"alerts"`
	LicenseDetails []any                  `json:"licenseDetails"`
}

type Vuln struct {
	// OSV analyzer
	ID                   string               `json:"id,omitempty"`
	Summary              string               `json:"summary,omitempty"`
	Details              string               `json:"details,omitempty"`
	Aliases              []string             `json:"aliases,omitempty"`
	Modified             time.Time            `json:"modified,omitempty"`
	Published            time.Time            `json:"published,omitempty"`
	Related              []string             `json:"related,omitempty"`
	GhsaDatabaseSpecific GhsaDatabaseSpecific `json:"database_specific,omitempty"`
	References           []References         `json:"references,omitempty"`
	Affected             []Affected           `json:"affected,omitempty"`
	SchemaVersion        string               `json:"schema_version,omitempty"`
	CvssSeverity         []CvssSeverity       `json:"severity,omitempty"`

	// TODO: create a common function instead of interface that'll handle multiple components

	// EPSS Score
	Epss Epss `json:"epss,omitempty"`
}
