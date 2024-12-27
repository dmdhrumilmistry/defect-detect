package main

import (
	"flag"
	"os"

	"github.com/dmdhrumilmistry/defect-detect/pkg/sbomconvert"
	"github.com/rs/zerolog/log"
)

func main() {
	inputFilePath := flag.String("f", "", "input sbom file path")
	outputFilePath := flag.String("o", "output.json", "output sbom file path")
	flag.Parse()

	log.Info().Str("input", *inputFilePath).Str("output", *outputFilePath).Msg("")

	input, err := os.Open(*inputFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read input file")
	}

	output, err := os.Create(*outputFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open output file")
	}

	err = sbomconvert.ConvertSbom(input, output)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to convert SBOM format")
	}
}
