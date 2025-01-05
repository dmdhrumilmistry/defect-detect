package epss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/dmdhrumilmistry/defect-detect/pkg/utils"
	"github.com/rs/zerolog/log"
)

type EpssAnalyzer struct {
	BaseUrl      string
	EpssEndpoint string
}

func NewEpssAnalyzer() *EpssAnalyzer {
	return &EpssAnalyzer{
		BaseUrl:      "https://api.first.org",
		EpssEndpoint: "/data/v1/epss",
	}
}

func (a *EpssAnalyzer) GetEpssFromVuln(cveId string) (epss types.Epss, err error) {
	log.Info().Msgf("Processing EPSS score for CVE %s", cveId)
	queryParams := url.Values{}
	queryParams.Add("cve", cveId)

	apiUrl := a.BaseUrl + a.EpssEndpoint + "?" + queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate request for fetching epss for cve: %s", cveId)
		return epss, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch EPSS")
		return epss, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected HTTP status: %d", res.StatusCode)
		log.Error().Err(err).Msg("failed to fetch EPSS")
		return epss, err
	}

	var apiResp types.EpssApiResponseSchema
	if err = json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("failed to decode EPSS response")
	}

	if len(apiResp.Data) > 0 {
		epss = apiResp.Data[0]
	}

	return epss, nil
}

func (a *EpssAnalyzer) WorkerEpss(ch <-chan *types.Vuln, resultCh chan error) {

	for vuln := range ch {
		var err error
		var epss types.Epss

		ids := []string{vuln.ID}
		ids = append(ids, vuln.Aliases...)

		cveId := utils.FindRegexMatchEle(`CVE-\d{4}-\d{4,}`, ids)
		if cveId != "" {
			epss, err = a.GetEpssFromVuln(cveId)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch epss score")
				resultCh <- err
			}
		}

		vuln.Epss = epss
		log.Info().Msgf("EPSS Score for CVE %s: %v", cveId, epss)

		resultCh <- nil
	}
}

// TODO: pass vulns slice by reference
func (a *EpssAnalyzer) ProcessEpssForVulns(vulns []types.Vuln, workers int) []types.Vuln {
	vulnChan := make(chan *types.Vuln)
	errChan := make(chan error)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.WorkerEpss(vulnChan, errChan)
		}()
	}

	// Send original slice items to channel
	go func() {
		for i := range vulns {
			vulnChan <- &vulns[i] // Pass pointer to original slice item
		}
		close(vulnChan)
	}()

	// Close errChan after workers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Process errors
	for err := range errChan {
		if err != nil {
			log.Error().Err(err).Msg("error occurred while fetching epss score")
		}
	}

	// Log updated vulnerabilities
	for _, vuln := range vulns {
		log.Debug().Msgf("%v", vuln.Epss)
	}

	return vulns
}
