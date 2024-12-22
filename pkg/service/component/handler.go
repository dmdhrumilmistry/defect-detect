package component

import (
	"net/http"
	"strconv"

	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComponentHandler struct {
	store     types.ComponentStore
	sbomStore types.SbomStore
}

func NewComponentHandler(store types.ComponentStore, sbomStore types.SbomStore) *ComponentHandler {
	return &ComponentHandler{
		store:     store,
		sbomStore: sbomStore,
	}
}

func (s *ComponentHandler) RegisterRoutes(r *gin.Engine) {
	// api v1
	r.POST("/api/v1/component", s.AddComponentUsingSbomId)
	r.GET("/api/v1/component", s.GetComponents)
	r.GET("/api/v1/component/:id", s.GetComponentById)
	r.GET("/api/v1/component/getByName", s.GetComponentByName)
	log.Info().Msg("component routes registered")
}

// curl -X POST -F "sbom=@example-sbom.json" http://localhost:8080/api/v1/component
func (s *ComponentHandler) AddComponentUsingSbomId(c *gin.Context) {
	sbomId, exists := c.GetQuery("sbom_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	sbom, err := s.sbomStore.GetSbomById(sbomId, 5)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	Ids, err := s.store.AddComponentUsingSbom(sbom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add components from sbom"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Components created successfully from Sbom", "ids": Ids})
}

// curl "http://localhost:8080/api/v1/component?page=1&limit=10"
func (s *ComponentHandler) GetComponents(c *gin.Context) {
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

	sboms, err := s.store.GetPaginatedComponents(page, limit, 5)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse sbom data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	total, err := s.store.GetComponentTotalCount()
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

// curl http://localhost:8080/api/v1/component/{component_id}
func (s *ComponentHandler) GetComponentById(c *gin.Context) {
	// Get the ID from the path parameter
	idParam := c.Param("id")

	// Convert the string ID to a MongoDB ObjectID
	components, err := s.store.GetComponentById(idParam, 5)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "component not found"})
		return
	} else if err != nil {
		log.Error().Err(err).Msg("failed to fetch component")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch component"})
		return
	}

	// Return the item as JSON
	c.JSON(http.StatusOK, components)
}

// curl "http://localhost:8080/api/v1/component/getByName?name=enigma"
func (s *ComponentHandler) GetComponentByName(c *gin.Context) {
	// Get the ID from the path parameter
	name, exists := c.GetQuery("name")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Convert the string ID to a MongoDB ObjectID
	components, err := s.store.GetComponentByName(name, 5)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	// Return the item as JSON
	c.JSON(http.StatusOK, components)
}