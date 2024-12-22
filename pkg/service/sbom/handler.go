package sbom

import (
	"net/http"
	"strconv"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComponentSbomHandler struct {
	store types.SbomStore
}

func NewComponentSbomHandler(store types.SbomStore) *ComponentSbomHandler {
	return &ComponentSbomHandler{
		store: store,
	}
}

func (s *ComponentSbomHandler) RegisterRoutes(r *gin.Engine) {
	// api v1
	r.POST("/api/v1/sbom", s.UploadSbomHandler)
	r.GET("/api/v1/sbom", s.GetSboms)
	r.GET("/api/v1/sbom/:id", s.GetSbomById)
	r.GET("/api/v1/sbom/getByComponentName", s.GetSbomByName)
	log.Info().Msg("sbom routes registered")
}

// curl -X POST -F "sbom=@example-sbom.json" http://localhost:8080/api/v1/sbom
func (s *ComponentSbomHandler) UploadSbomHandler(c *gin.Context) {
	file, err := c.FormFile("sbom")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file content"})
		return
	}
	defer fileContent.Close()

	decoder := cyclonedx.NewBOMDecoder(fileContent, cyclonedx.BOMFileFormatJSON)
	var bom cyclonedx.BOM
	if err := decoder.Decode(&bom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid SBOM format"})
		return
	}

	// Store component SBOM
	componentId, err := s.store.AddComponentSbom(bom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upload component SBOM",
		})
		return
	}
	// TODO: create scan task

	c.JSON(http.StatusOK, gin.H{"message": "SBOM uploaded successfully", "id": componentId})
}

// curl http://localhost:8080/api/v1/sbom
func (s *ComponentSbomHandler) GetSboms(c *gin.Context) {
	// Get page and limit from query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit >= 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	sboms, err := s.store.GetPaginatedSboms(page, limit, 5)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse sbom data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	total, err := s.store.GetComponentSbomTotalCount()
	if err != nil {
		log.Error().Err(err).Msg("failed to get total sbom data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	// Build response
	c.JSON(http.StatusOK, gin.H{
		"data":  sboms,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// curl http://localhost:8080/api/v1/sbom/{sbom_id}
func (s *ComponentSbomHandler) GetSbomById(c *gin.Context) {
	// Get the ID from the path parameter
	idParam := c.Param("id")

	// Convert the string ID to a MongoDB ObjectID
	sbom, err := s.store.GetSbomById(idParam, 5)
	log.Print(err)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	// Return the item as JSON
	c.JSON(http.StatusOK, sbom)
}

// curl "http://localhost:8080/api/v1/sbom/getByComponentName?name=enigma"
func (s *ComponentSbomHandler) GetSbomByName(c *gin.Context) {
	// Get the ID from the path parameter
	name, exists := c.GetQuery("name")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Convert the string ID to a MongoDB ObjectID
	sboms, err := s.store.GetSbomByName(name, 5)
	log.Print(err)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	// Return the item as JSON
	c.JSON(http.StatusOK, sboms)
}
