package main

import (
	"context"

	"github.com/dmdhrumilmistry/defect-detect/pkg/analyzer/osv"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/db"
	"github.com/dmdhrumilmistry/defect-detect/pkg/service/component"
	"github.com/dmdhrumilmistry/defect-detect/pkg/service/sbom"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r = r.With()
	r.SetTrustedProxies(nil)

	mgo, err := db.NewMongo(config.DefaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get db connection")
	}
	defer mgo.Client.Disconnect(context.TODO())

	if !config.DefaultConfig.IsDevEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Analyzers
	osvAnalyzer := osv.NewOsvAnalyzer()

	// create stores
	log.Info().Msg("Registering Routes")
	sbomStore := sbom.NewComponentSbomStore(mgo.Db)
	sbomHandler := sbom.NewComponentSbomHandler(sbomStore)
	sbomHandler.RegisterRoutes(r)

	componentStore := component.NewComponentStore(mgo.Db, osvAnalyzer)
	componentHandler := component.NewComponentHandler(componentStore, sbomStore)
	componentHandler.RegisterRoutes(r)

	// Start the server
	if err := r.Run(":" + config.DefaultConfig.HostPort); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start server on port %s", config.DefaultConfig.HostPort)
	}
}
