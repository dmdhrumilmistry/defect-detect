package osv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
)

type OsvAnalyzer struct {
	baseUrl string
}

func NewOsvAnalyzer() *OsvAnalyzer {
	return &OsvAnalyzer{
		baseUrl: "https://api.osv.dev",
	}
}

func (a *OsvAnalyzer) GetVulns(purl string) ([]types.Vuln, error) {
	vulns := []types.Vuln{}
	resp, err := a.getVuln(purl, "")
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch vulns for purl: %s", purl)
		return vulns, err
	}

	vulns = append(vulns, resp.Vulns...)

	for resp.NextPageToken != "" {
		resp, err = a.getVuln(purl, resp.NextPageToken)
		if err != nil {
			log.Error().Err(err).Msgf("failed to fetch vulns for purl: %s", purl)
			continue
		}
		vulns = append(vulns, resp.Vulns...)
	}

	return vulns, nil

}

// queries OSV api and fetches
func (a *OsvAnalyzer) getVuln(purl string, pageToken string) (types.OsvQueryApiResponse, error) {
	log.Info().Msgf("Fetching vulns for purl %s with page token %s", purl, pageToken)
	osvResp := types.OsvQueryApiResponse{}
	url := a.baseUrl + "/v1/query"

	payload := map[string]interface{}{
		"package": map[string]string{
			"purl": purl,
		},
	}

	if pageToken != "" {
		payload["page_token"] = pageToken
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal json payload data for purl %s with page token %s", purl, pageToken)
		return osvResp, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch vuln data for purl %s with page token %s", purl, pageToken)
		return osvResp, err
	}

	if response.StatusCode != http.StatusOK {
		log.Error().Err(err).Msgf("failed to fetch vuln data for purl %s with page token %s. OSV api returned status code %d instead of 200", purl, pageToken, response.StatusCode)
		return osvResp, fmt.Errorf("OSV api returned status code %d instead of 200", response.StatusCode)
	}
	defer response.Body.Close()

	return osvResp, json.NewDecoder(response.Body).Decode(&osvResp)
}
