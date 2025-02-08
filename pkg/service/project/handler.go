package project

import (
	"fmt"
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
	store          types.ProjectStore
	componentStore types.ComponentStore
	sbomStore      types.SbomStore
	authStore      types.AuthStore
}

func NewProjectHandler(store types.ProjectStore, sbomStore types.SbomStore, componentStore types.ComponentStore, authStore types.AuthStore) *ProjectHandler {
	return &ProjectHandler{
		store:          store,
		componentStore: componentStore,
		sbomStore:      sbomStore,
		authStore:      authStore,
	}
}

func (p *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	// api v1
	r.POST("/api/v1/project", p.CreateProject)
	r.GET("/api/v1/project", p.GetProjects)
	r.GET("/api/v1/project/:id", p.GetProjectById)
	r.PATCH("/api/v1/project/:id", p.UpdateProjectById)
	r.DELETE("/api/v1/project/:id", p.DeleteProjectById)

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
func (p *ProjectHandler) GetProjectById(c *gin.Context) {
	// Get the ID from the path parameter
	idParam := c.Param("id")

	// Convert the string ID to a MongoDB ObjectID
	projects, err := p.store.GetProjectById(idParam, config.DefaultConfig.DbQueryTimeout)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	} else if err != nil {
		log.Error().Err(err).Msg("failed to fetch project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	// Return the item as JSON
	if len(projects) > 0 {
		c.JSON(http.StatusOK, projects[0])
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err": "project not found",
	})
}

func (p *ProjectHandler) UpdateProjectById(c *gin.Context) {
	var payload types.Project
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate payload"})
		return
	}

	idParam := c.Param("id")
	payload.Id = idParam
	if err := p.store.UpdateById(payload, config.DefaultConfig.DbQueryTimeout); err != nil {
		log.Error().Err(err).Msgf("failed to update project details for id %s: %v", idParam, payload)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update project details",
		})
		return
	}

	if len(payload.Sboms) > payload.SbomsToRetain {
		// latest sboms should be present at 0th position
		payload.Sboms = payload.Sboms[:payload.SbomsToRetain]
	}

	if err := p.sbomStore.ValidateIds(payload.Sboms); err != nil {
		log.Error().Err(err).Msgf("failed to validate sbom ids: %v", payload.Sboms)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate sbom ids"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "project details updated successfully"})
}

// curl http://localhost:8080/api/v1/project/{project_id}
func (s *ProjectHandler) DeleteProjectById(c *gin.Context) {
	// Get the ID from the path parameter
	var deleteSboms bool
	var msg string
	idParam := c.Param("id")
	deleteSbom, _ := c.GetQuery("delete_sbom")

	deleteSboms, err := strconv.ParseBool(deleteSbom)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse delete_sbom param. Using default value: false")
		deleteSboms = false
	}

	if deleteSboms {
		projects, err := s.store.GetProjectById(idParam, config.DefaultConfig.DbQueryTimeout)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		} else if err != nil {
			log.Error().Err(err).Msg("failed to fetch project")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
			return
		}

		deleteCount, err := s.componentStore.DeleteByIds(projects[0].Sboms, "sbom_id", config.DefaultConfig.DbQueryTimeout)
		if err != nil {
			log.Error().Err(err).Msg("failed to delete sboms components")
			msg += "failed to delete sbom components."
		} else {
			msg += fmt.Sprintf("%d sbom components deleted successfully.", deleteCount)
			log.Info().Msgf("Total %d sbom components deleted", deleteCount)
		}

		deleteCount, err = s.sbomStore.DeleteByIds(projects[0].Sboms, config.DefaultConfig.DbQueryTimeout)
		if err != nil {
			log.Error().Err(err).Msg("failed to delete sboms")
			msg += "failed to delete sboms."
		} else {
			msg += fmt.Sprintf("%d sboms deleted successfully.", deleteCount)
			log.Info().Msgf("Total %d sboms deleted", deleteCount)
		}

	}

	_, err = s.store.DeleteById(idParam, config.DefaultConfig.DbQueryTimeout)
	if err != nil {
		log.Error().Err(err).Msgf("failed to delete project with id: %s", idParam)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete project. %s", msg)})
		return
	}
	log.Info().Msgf("Project %s deleted successfully", idParam)

	c.JSON(http.StatusNoContent, gin.H{
		"msg": fmt.Sprintf("Project deleted successfully. %s", msg),
	})

}
