package types

type Analyzer interface {
	GetVulns(purl string) ([]Vuln, error)
}

// Auto generated struct code for OSV response schema
type OsvQueryApiResponse struct {
	Vulns         []Vuln `json:"vulns,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}
type References struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Package struct {
	Name      string `json:"name,omitempty"`
	Ecosystem string `json:"ecosystem,omitempty"`
	Purl      string `json:"purl,omitempty"`
}
type Events struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}
type Ranges struct {
	Type   string   `json:"type,omitempty"`
	Repo   string   `json:"repo,omitempty"`
	Events []Events `json:"events,omitempty"`
}
type EcosystemSpecific struct {
	Severity string `json:"severity,omitempty"`
}
type DatabaseSpecific struct {
	Source string `json:"source,omitempty"`
}
type Affected struct {
	Package           Package           `json:"package,omitempty"`
	Ranges            []Ranges          `json:"ranges,omitempty"`
	Versions          []string          `json:"versions,omitempty"`
	EcosystemSpecific EcosystemSpecific `json:"ecosystem_specific,omitempty"`
	DatabaseSpecific  DatabaseSpecific  `json:"database_specific,omitempty"`
}
