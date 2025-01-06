package mpaf

import (
	"slices"

	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/dmdhrumilmistry/m-paf/pkg/socketdev"
	"github.com/rs/zerolog/log"
)

type MpafAnalyzer struct {
	Api *socketdev.Api
}

func NewMpafAnalyzer() (*MpafAnalyzer, error) {
	api, err := socketdev.NewSocketAPI()
	if err != nil {
		log.Error().Err(err).Msg("failed to init socket api")
		return nil, err
	}

	return &MpafAnalyzer{
		Api: api,
	}, nil
}

func (a *MpafAnalyzer) mapAlerts(alerts []socketdev.Alert) []socketdev.AlertType {
	var alertTypes []socketdev.AlertType
	var alertTypeIds []int

	for _, alert := range alerts {
		alertType, ok := a.Api.AlertTypes[alert.Type]
		if !ok {
			log.Error().Msgf("Alert Type with id %d not found", alert.Type)
		}

		if !slices.Contains(alertTypeIds, alert.Type) {
			alertTypeIds = append(alertTypeIds, alert.Type)
			alertTypes = append(alertTypes, alertType)
		}

	}

	return alertTypes
}

func (a *MpafAnalyzer) GetPackageInfo(purl string) ([]types.PackageInfo, error) {
	var pkgInfos []types.PackageInfo
	interPkgInfos, err := a.Api.GetAlerts(purl)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get packageInfo for purl: %s", purl)
		return pkgInfos, err
	}

	for _, pkgInfo := range interPkgInfos {
		pkgInfos = append(pkgInfos, types.PackageInfo{
			// ID:             pkgInfo.ID,
			// Type:           pkgInfo.Type,
			// Name:           pkgInfo.Name,
			// Namespace:      pkgInfo.Namespace,
			// Files:          pkgInfo.Files,
			// Version:        pkgInfo.Version,
			// Qualifiers:     pkgInfo.Qualifiers,
			Scores:         pkgInfo.Scores,
			Capabilities:   pkgInfo.Capabilities,
			License:        pkgInfo.License,
			Size:           pkgInfo.Size,
			State:          pkgInfo.State,
			Alerts:         a.mapAlerts(pkgInfo.Alerts),
			LicenseDetails: pkgInfo.LicenseDetails,
		})
	}

	return pkgInfos, nil
}
