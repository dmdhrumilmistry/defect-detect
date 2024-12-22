package main

import (
	"context"
	"net/http"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/db"
	"github.com/dmdhrumilmistry/defect-detect/pkg/service/component"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SBOMData struct {
	Components []cyclonedx.Component `json:"components"`
}

var sboms = make(map[string]SBOMData) // In-memory storage for simplicity

// curl "http://localhost:8080/api/components?sbom_id=bom.json"
func listComponents(c *gin.Context) {
	sbomID := c.Query("sbom_id")
	sbomData, exists := sboms[sbomID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "SBOM not found"})
		return
	}

	c.JSON(http.StatusOK, sbomData.Components)
}

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r = r.With()
	r.SetTrustedProxies(nil)

	cfg := config.NewConfig()
	mgo, err := db.NewMongo(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get db connection")
	}
	defer mgo.Client.Disconnect(context.TODO())

	if !cfg.IsDevEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// create stores
	compStore := component.NewComponentSbomStore(mgo.Db)
	compHandler := component.NewComponentSbomHandler(compStore)
	compHandler.RegisterRoutes(r)

	// Start the server
	if err := r.Run(":" + cfg.HostPort); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start server on port %s", cfg.HostPort)
	}
}
