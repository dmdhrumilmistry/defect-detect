package types

import (
	"bytes"
	"encoding/xml"

	"github.com/CycloneDX/cyclonedx-go"
)

type SbomStore interface {
	AddComponentSbom(sbom cyclonedx.BOM) (string, error)
	GetComponentSbomTotalCount() (int64, error)
	GetPaginatedSboms(page, limit, duration int) ([]Sbom, error)
	GetSbomById(idParam string, duration int) (Sbom, error)
	GetSbomByName(name string, duration int) ([]Sbom, error)
}

type Sbom struct {
	Id string `json:"sbom_id" bson:"_id"`

	// XML specific fields
	XMLName xml.Name `json:"-" xml:"bom"`
	XMLNS   string   `json:"-" xml:"xmlns,attr"`

	// JSON specific fields
	JSONSchema  string                `json:"$schema,omitempty" xml:"-"`
	BOMFormat   string                `json:"bomFormat" xml:"-"`
	SpecVersion cyclonedx.SpecVersion `json:"specVersion" xml:"-"`

	SerialNumber       string                         `json:"serialNumber,omitempty" xml:"serialNumber,attr,omitempty"`
	Version            int                            `json:"version" xml:"version,attr"`
	Metadata           *cyclonedx.Metadata            `json:"metadata,omitempty" xml:"metadata,omitempty"`
	Components         *[]cyclonedx.Component         `json:"components,omitempty" xml:"components>component,omitempty"`
	Services           *[]cyclonedx.Service           `json:"services,omitempty" xml:"services>service,omitempty"`
	ExternalReferences *[]cyclonedx.ExternalReference `json:"externalReferences,omitempty" xml:"externalReferences>reference,omitempty"`
	Dependencies       *[]cyclonedx.Dependency        `json:"dependencies,omitempty" xml:"dependencies>dependency,omitempty"`
	Compositions       *[]cyclonedx.Composition       `json:"compositions,omitempty" xml:"compositions>composition,omitempty"`
	Properties         *[]cyclonedx.Property          `json:"properties,omitempty" xml:"properties>property,omitempty"`
	Vulnerabilities    *[]cyclonedx.Vulnerability     `json:"vulnerabilities,omitempty" xml:"vulnerabilities>vulnerability,omitempty"`
	Annotations        *[]cyclonedx.Annotation        `json:"annotations,omitempty" xml:"annotations>annotation,omitempty"`
	Formulation        *[]cyclonedx.Formula           `json:"formulation,omitempty" xml:"formulation>formula,omitempty"`
	Declarations       *cyclonedx.Declarations        `json:"declarations,omitempty" xml:"declarations,omitempty"`
	Definitions        *cyclonedx.Definitions         `json:"definitions,omitempty" xml:"definitions,omitempty"`
}

type GithubRepoImportRequestSchema struct {
	Owner    string `json:"owner" binding:"required"`
	RepoName string `json:"repo_name" binding:"required"`
}

type GithubRepoImportResponseSchema struct {
	Sbom map[string]interface{}
}

type ReadSeekCloser struct {
	*bytes.Reader
}

// Implement the Close method (no-op)
func (r *ReadSeekCloser) Close() error {
	// No-op: bytes.Reader doesn't need to be closed
	return nil
}

// Custom WriteCloser type to implement io.WriteCloser
type WriteCloser struct {
	*bytes.Buffer
}

// Implement the Close method (no-op for Buffer)
func (w *WriteCloser) Close() error {
	// No-op: bytes.Buffer doesn't need to be closed
	return nil
}
