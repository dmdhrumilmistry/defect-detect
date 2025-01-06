package types

import "time"

type Analyzer interface {
	GetVulns(purl string) ([]Vuln, error)
	GetPackageInfo(purl string) ([]PackageInfo, error)
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

type GhsaDatabaseSpecific struct {
	GithubReviewedAt time.Time `json:"github_reviewed_at,omitempty"`
	GithubReviewed   bool      `json:"github_reviewed,omitempty"`
	Severity         string    `json:"severity,omitempty"`
	CweIds           []string  `json:"cwe_ids,omitempty"`
	NvdPublishedAt   time.Time `json:"nvd_published_at,omitempty"`
}

type CvssSeverity struct {
	TypeStr  string `json:"type,omitempty"`
	ScoreStr string `json:"score,omitempty"`
}

type Affected struct {
	Package           Package           `json:"package,omitempty"`
	Ranges            []Ranges          `json:"ranges,omitempty"`
	Versions          []string          `json:"versions,omitempty"`
	EcosystemSpecific EcosystemSpecific `json:"ecosystem_specific,omitempty"`
	DatabaseSpecific  DatabaseSpecific  `json:"database_specific,omitempty"`
}

// End of Auto generated struct code for OSV response schema

// Start of EPSS Structs
type Epss struct {
	CveId      string `json:"cve,omitempty" bson:"cve,omitempty"`
	EpssScore  string `json:"epss,omitempty" bson:"epss,omitempty"`
	Percentile string `json:"percentile,omitempty" bson:"percentile,omitempty"`
	Date       string `json:"date,omitempty" bson:"date,omitempty"`
}

type EpssApiResponseSchema struct {
	// Status                    string `json:"status,omitempty"`
	StatusCode int `json:"status-code,omitempty"`
	// Version                   string `json:"version,omitempty"`
	AccessControlAllowHeaders string `json:"access-control-allow-headers,omitempty"`
	// Access                    string `json:"access,omitempty"`
	Total  int    `json:"total,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Data   []Epss `json:"data,omitempty"`
}

// End of EPSS Structs
