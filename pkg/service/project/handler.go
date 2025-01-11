package project

import (
	"net/http"
	"strconv"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectHandler struct {
	store *ProjectStore
}

func NewProjectHandler(store *ProjectStore) *ProjectHandler {
	return &ProjectHandler{
		store: store,
	}
}

func (p *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	// api v1
	r.POST("/api/v1/project", p.CreateProject)
	r.GET("/api/v1/project", p.GetProjects)
	r.GET("/api/v1/project/:id", p.GetProjectById)
	// TODO: implement routes for updating and deleting project

	log.Info().Msg("Project routes registered")
}

func (p *ProjectHandler) CreateProject(c *gin.Context) {
	var payload types.Project
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate payload"})
		return
	}

	// ignore provided id
	payload.Id = ""

	id, err := p.store.AddProject(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create project",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project created successfully", "id": id})
}

// curl http://localhost:8080/api/v1/project
func (p *ProjectHandler) GetProjects(c *gin.Context) {
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

	filter := bson.M{}
	sboms, err := p.store.GetUsingFilter(filter, page, limit, config.DefaultConfig.DbQueryTimeout)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse project data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	total, err := p.store.GetTotalCount(filter)
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

// curl http://localhost:8080/api/v1/project/{project_id}
func (s *ProjectHandler) GetProjectById(c *gin.Context) {
	// Get the ID from the path parameter
	idParam := c.Param("id")

	// Convert the string ID to a MongoDB ObjectID
	projects, err := s.store.GetProjectById(idParam, config.DefaultConfig.DbQueryTimeout)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	} else if err != nil {
		log.Error().Err(err).Msg("failed to fetch project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Return the item as JSON
	c.JSON(http.StatusOK, projects)
}
