package analyzer

import (
	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/epss"
	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/mpaf"
	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/osv"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
)

type Analyzer struct {
	RunOsv  bool
	RunMpaf bool
	RunEpss bool

	OsvAnalyzer  *osv.OsvAnalyzer
	MpafAnalyzer *mpaf.MpafAnalyzer
	EpssAnalyzer *epss.EpssAnalyzer
}

func NewAnalyzer() *Analyzer {
	mpafAnalyzer, err := mpaf.NewMpafAnalyzer()
	if err != nil {
		log.Error().Err(err).Msg("failed to init mpaf analyzer")
	}

	return &Analyzer{
		// config
		RunOsv:  config.DefaultConfig.RunOsv,
		RunMpaf: config.DefaultConfig.RunMpaf,
		RunEpss: config.DefaultConfig.RunEpss,

		// Analyzers
		OsvAnalyzer:  osv.NewOsvAnalyzer(),
		MpafAnalyzer: mpafAnalyzer,
		EpssAnalyzer: epss.NewEpssAnalyzer(),
	}
}

func (a *Analyzer) GetPackageInfo(purl string) (pkgInfos []types.PackageInfo, err error) {
	if a.RunMpaf {
		log.Info().Msgf("Fetching package info for purl: %s", purl)
		return a.MpafAnalyzer.GetPackageInfo(purl)
	} else {
		log.Error().Msgf("mpaf analyzer is not enabled. Skipping fetching package info for purl: %s", purl)
	}

	return pkgInfos, nil
}

func (a *Analyzer) GetVulns(purl string) (vulns []types.Vuln, err error) {
	log.Info().Msgf("Running analyzers for purl: %s", purl)
	if a.RunOsv {
		vulns, err = a.OsvAnalyzer.GetVulns(purl)
		if err != nil {
			log.Error().Err(err).Msgf("failed to retrieve osv vulns for purl: %s", purl)
		}
	}

	// concurrently update epss for cvss
	if a.RunEpss && len(vulns) > 0 {
		log.Info().Msgf("running epss analyzer on vulns for purl: %s", purl)
		vulns = a.EpssAnalyzer.ProcessEpssForVulns(vulns, config.DefaultConfig.DefaultWorkersCount)
	}

	log.Info().Msgf("Completed analysis for purl: %s", purl)

	return vulns, nil
}
