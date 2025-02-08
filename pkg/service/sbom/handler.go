package sbom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/sbomconvert"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComponentSbomHandler struct {
	authStore types.AuthStore
	store     types.SbomStore
}

func NewComponentSbomHandler(store types.SbomStore, authStore types.AuthStore) *ComponentSbomHandler {
	return &ComponentSbomHandler{
		store:     store,
		authStore: authStore,
	}
}

func (s *ComponentSbomHandler) RegisterRoutes(r *gin.Engine, authStore types.AuthStore) {
	// api v1
	r.POST("/api/v1/sbom", s.UploadSbomHandler)
	r.GET("/api/v1/sbom", s.GetSboms)
	r.GET("/api/v1/sbom/:id", s.GetSbomById)
	r.GET("/api/v1/sbom/getByComponentName", s.GetSbomByName)
	r.POST("/api/v1/sbom/convert", s.ConvertSbom)
	r.POST("/api/v1/sbom/githubImport", s.ImportGithubRepo)

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

	sboms, err := s.store.GetPaginatedSboms(page, limit, config.DefaultConfig.DbQueryTimeout)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse sbom data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	total, err := s.store.GetTotalCount(bson.M{})
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
	sbom, err := s.store.GetSbomById(idParam, config.DefaultConfig.DbQueryTimeout)
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
	sboms, err := s.store.GetSbomByName(name, config.DefaultConfig.DbQueryTimeout)
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

// curl -X POST -F "sbom=@sbom.json" "http://localhost:8080/api/v1/sbom/convert"
// curl -X POST -F "sbom=@sbom.json" "http://localhost:8080/api/v1/sbom/convert" | jq '.converted_sbom.components.[].purl' -r
func (s *ComponentSbomHandler) ConvertSbom(c *gin.Context) {
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

	// Create a buffer to store the converted SBOM output
	// Create a pipe to use as the output stream with a Close method
	outputBuffer := &bytes.Buffer{}
	outputStream := struct {
		io.Writer
		io.Closer
	}{
		Writer: outputBuffer,
		Closer: io.NopCloser(nil), // Provides a no-op Close method
	}

	// Perform the conversion
	err = sbomconvert.ConvertSbom(fileContent, outputStream)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert sbom")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert SBOM"})
		return
	}

	var sbom map[string]interface{}
	if err := json.Unmarshal(outputBuffer.Bytes(), &sbom); err != nil {
		log.Error().Err(err).Msg("failed to convert sbom")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert SBOM"})
		return
	}

	// update uuid for spdx
	_, exists := sbom["serialNumber"]
	if exists {
		sbom["serialNumber"] = fmt.Sprintf("urn:uuid:%s", uuid.New().String())
	}

	// Return the converted SBOM as JSON
	c.JSON(http.StatusOK, gin.H{"converted_sbom": sbom})
}

// curl -X POST -H "application/json" -d '{"owner":"dmdhrumilmistry", "repo_name":"pyhtools"}' http://localhost:8080/api/v1/sbom/githubImport
func (s *ComponentSbomHandler) ImportGithubRepo(c *gin.Context) {
	var jsonData types.GithubRepoImportRequestSchema

	// Bind the incoming JSON data to the struct
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// If there is an error parsing the JSON, send an error response
		log.Error().Err(err).Msg("invalid json data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate the URL
	sbomUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/dependency-graph/sbom", jsonData.Owner, jsonData.RepoName)
	regex := `^https:\/\/api\.github\.com\/repos\/[a-zA-Z0-9_-]+\/[a-zA-Z0-9_-]+\/dependency-graph\/sbom$`
	re := regexp.MustCompile(regex)

	if !re.MatchString(sbomUrl) {
		log.Error().Msgf("owner (%s) and repo name (%s) are invalid", jsonData.Owner, jsonData.RepoName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	log.Info().Msgf("Fetching SBOM for https://github.com/%s/%s repo", jsonData.Owner, jsonData.RepoName)

	req, err := http.NewRequest(http.MethodGet, sbomUrl, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate req for github.com/%s/%s", jsonData.Owner, jsonData.RepoName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate request to fetch sbom github api"})
		return
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.DefaultConfig.GithubToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch sbom from github.com/%s/%s", jsonData.Owner, jsonData.RepoName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sbom github api"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Err(err).Msgf("failed to fetch sbom from github.com/%s/%s. Expected 200 received %d from api", jsonData.Owner, jsonData.RepoName, res.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sbom github api"})
		return
	}

	// Read the response body into a buffer
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read github api response body into buffer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the response body"})
		return
	}

	var respPayload types.GithubRepoImportResponseSchema
	if err := json.Unmarshal(body, &respPayload); err != nil {
		log.Error().Err(err).Msg("failed to parse github api resp to json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the response body"})
		return
	}

	sbomByte, err := json.Marshal(respPayload.Sbom)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert github api sbom resp to json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the github repo sbom"})

		return
	}

	// Convert the buffer into an io.ReadSeekCloser (using a bytes.Reader)
	readSeekCloser := &types.ReadSeekCloser{Reader: bytes.NewReader(sbomByte)}

	// Create a buffer to store the converted SBOM output
	outputBuffer := &bytes.Buffer{}
	outputStream := &types.WriteCloser{
		Buffer: outputBuffer,
	}

	if err := sbomconvert.ConvertSbom(readSeekCloser, outputStream); err != nil {
		log.Error().Err(err).Msg("failed to convert sbom")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert sbom"})
		return
	}

	decoder := cyclonedx.NewBOMDecoder(outputStream, cyclonedx.BOMFileFormatJSON)
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

	// TODO: add task to process sbom component

	c.JSON(http.StatusOK, gin.H{"message": "SBOM uploaded successfully", "id": componentId})
}
