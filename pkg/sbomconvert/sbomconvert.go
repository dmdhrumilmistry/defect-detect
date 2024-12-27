package sbomconvert

import (
	"context"
	"io"

	"github.com/protobom/sbom-convert/pkg/convert"
	"github.com/protobom/sbom-convert/pkg/format"
	"github.com/rs/zerolog/log"
)

func ConvertSbom(inputStream io.ReadSeekCloser, outputStream io.WriteCloser) error {
	frmt, err := format.Detect(inputStream)
	if err != nil {
		log.Error().Err(err).Msg("Failed to detect format")
		return err
	}
	log.Info().Msgf("SBOM format detected %s", frmt.String())

	inverse, err := frmt.Inverse()
	if err != nil {
		log.Error().Err(err).Msg("failed while performing inverse")
	}
	log.Info().Msgf("Converting SBOM to %s", inverse.String())

	service := convert.NewService(
		convert.WithFormat(inverse),
	)

	return service.Convert(context.TODO(), inputStream, outputStream)
}
