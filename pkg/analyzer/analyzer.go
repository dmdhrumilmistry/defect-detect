package analyzer

import (
	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/mpaf"
	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/osv"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
)

type Analyzer struct {
	RunOsv  bool
	RunMpaf bool

	OsvAnalyzer  *osv.OsvAnalyzer
	MpafAnalyzer *mpaf.MpafAnalyzer
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

		// Analyzers
		OsvAnalyzer:  osv.NewOsvAnalyzer(),
		MpafAnalyzer: mpafAnalyzer,
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
	if a.RunOsv {
		vulns, err = a.OsvAnalyzer.GetVulns(purl)
		if err != nil {
			log.Error().Err(err).Msgf("failed to retrieve osv vulns for purl: %s", purl)

		}
	}

	// TODO: concurrently update epss for cvss

	return vulns, nil
}
