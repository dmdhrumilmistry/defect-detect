package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SBOMData struct {
	Components []cyclonedx.Component `json:"components"`
}

var sboms = make(map[string]SBOMData) // In-memory storage for simplicity

// curl -X POST -F "sbom=@example-sbom.json" http://localhost:8080/api/sbom/upload
func uploadSBOM(c *gin.Context) {
	file, err := c.FormFile("sbom")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file content"})
	}
	defer fileContent.Close()

	decoder := cyclonedx.NewBOMDecoder(fileContent, cyclonedx.BOMFileFormatJSON)
	var bom cyclonedx.BOM
	if err := decoder.Decode(&bom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid SBOM format"})
	}
	fmt.Print(bom)

	// Store components in memory (replace with DB in production)
	sboms[file.Filename] = SBOMData{Components: *bom.Components}

	c.JSON(http.StatusOK, gin.H{"message": "SBOM uploaded successfully", "file": file.Filename})
}

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
	r := gin.Default()
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

	// Define API endpoints
	r.POST("/api/sbom/upload", uploadSBOM)
	r.GET("/api/components", listComponents)

	// Start the server
	if err := r.Run(":" + cfg.HostPort); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start server on port %s", cfg.HostPort)
	}
}
