package component

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/dmdhrumilmistry/defect-detect/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COMPONENT_COLLECTION = "component"

type vulnResult struct {
	Component types.Component
	Err       error
}
type ComponentStore struct {
	db         *mongo.Database
	collection *mongo.Collection
	Analyzer   types.Analyzer
}

func NewComponentStore(db *mongo.Database, analyzer types.Analyzer) *ComponentStore {
	collection := db.Collection(COMPONENT_COLLECTION)
	// TODO: create index if not exists

	return &ComponentStore{
		db:         db,
		collection: collection,
		Analyzer:   analyzer,
	}
}

func (c *ComponentStore) processComponentsWorker(sbom types.Sbom, componentName, componentVersion string, wg *sync.WaitGroup, workCh <-chan *cyclonedx.Component, resultCh chan vulnResult) {
	defer wg.Done()
	for component := range workCh {
		var licences []string
		var vulns []types.Vuln
		var pkgInfos []types.PackageInfo
		var vulnErr, pkgInfoErr error

		if component.Licenses != nil {
			for _, license := range *component.Licenses {
				if license.License != nil && license.License.ID != "" {
					licences = append(licences, license.License.ID)
				}
			}
		}

		var innerWg sync.WaitGroup
		// Channel to collect errors
		errCh := make(chan error, 2)

		// Fetch vulnerabilities concurrently
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
			if component.PackageURL != "" {
				log.Info().Msgf("Processing vulns for purl %s", component.PackageURL)
				vulns, vulnErr = c.Analyzer.GetVulns(component.PackageURL)
				if vulnErr != nil {
					log.Error().Err(vulnErr).Msgf("failed to analyze vulns for %s", component.PackageURL)
					errCh <- vulnErr
				} else {
					log.Info().Msgf("Detected %d vulns for purl: %s", len(vulns), component.PackageURL)
				}
			}
		}()

		// Fetch package information concurrently
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
			pkgInfos, pkgInfoErr = c.Analyzer.GetPackageInfo(component.PackageURL)
			if pkgInfoErr != nil {
				log.Error().Err(pkgInfoErr).Msgf("failed to fetch package info for purl: %s", component.PackageURL)
				errCh <- pkgInfoErr
			}
		}()

		// Wait for both goroutines to complete
		innerWg.Wait()
		close(errCh)

		// Aggregate errors
		var combinedErr error
		for err := range errCh {
			if combinedErr == nil {
				combinedErr = err
			} else {
				combinedErr = fmt.Errorf("%v; %w", combinedErr, err)
			}
		}

		// Send the result back
		resultCh <- vulnResult{
			Component: types.Component{
				Name:             component.Name,
				Version:          component.Version,
				PackageUrl:       component.PackageURL,
				Licenses:         licences,
				Type:             string(component.Type),
				ComponentName:    componentName,
				ComponentVersion: componentVersion,
				Vulns:            vulns,
				SbomId:           sbom.Id,
				PackageInfos:     pkgInfos,
			},
			Err: combinedErr,
		}
	}
}

func (c *ComponentStore) processComponents(sbom types.Sbom, componentName, componentVersion string, workers int) []interface{} {
	var components []interface{}

	// Channels for work distribution and results collection
	workCh := make(chan *cyclonedx.Component)
	resultCh := make(chan vulnResult)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		// go worker(&wg)
		go c.processComponentsWorker(sbom, componentName, componentVersion, &wg, workCh, resultCh)
	}

	// Send components to work channel
	go func() {
		for _, component := range *sbom.Components {
			workCh <- &component
		}
		close(workCh)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		components = append(components, result.Component)
	}

	return components
}

func (c *ComponentStore) IsSbomProcessed(sbomId string) bool {
	doc_count, err := c.collection.CountDocuments(context.TODO(), bson.M{"sbom_id": sbomId})
	if err != nil {
		log.Error().Err(err).Msgf("failed to count component docs with sbom_id %s", sbomId)
		return false
	}

	log.Info().Msgf("%d components are already processed %s", doc_count, sbomId)
	return doc_count != 0
}

func (c *ComponentStore) AddComponentUsingSbom(sbom types.Sbom) ([]string, error) {
	componentName := sbom.Metadata.Component.Name
	componentVersion := sbom.Metadata.Component.Version
	insertedIds := []string{}

	if c.IsSbomProcessed(sbom.Id) {
		return insertedIds, fmt.Errorf("sbom is already processed")
	}

	components := c.processComponents(sbom, componentName, componentVersion, config.DefaultConfig.DefaultWorkersCount)

	results, err := c.collection.InsertMany(context.TODO(), components)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert")
		return insertedIds, err
	}

	for _, insertedId := range results.InsertedIDs {
		insertedIds = append(insertedIds, insertedId.(primitive.ObjectID).Hex())
	}

	return insertedIds, nil
}

func (c *ComponentStore) GetComponentTotalCount(filter interface{}) (int64, error) {
	// Get total count of documents
	total, err := c.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return -1, err
	}

	return total, nil
}

func (c *ComponentStore) GetComponentsUsingFilter(filter interface{}, page, limit, duration int) ([]types.Component, error) {
	var components []types.Component

	// Calculate skip
	skip := (page - 1) * limit

	// MongoDB query options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	cursor, err := c.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return components, err
	}
	defer cursor.Close(ctx)

	// Parse results
	if err := cursor.All(ctx, &components); err != nil {
		return components, err
	}

	return components, nil
}

// Handler for getting paginated items
func (c *ComponentStore) GetPaginatedComponents(page, limit, duration int) ([]types.Component, error) {
	return c.GetComponentsUsingFilter(bson.M{}, page, limit, duration)
}

// Handler for getting paginated items
func (c *ComponentStore) GetComponentById(idParam string, duration int) ([]types.Component, error) {
	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return []types.Component{}, err
	}

	return c.GetComponentsUsingFilter(bson.M{"_id": objID}, 1, 1, config.DefaultConfig.DbQueryTimeout)
}

func (c *ComponentStore) GetComponentByName(name string, duration int) ([]types.Component, error) {
	return c.GetComponentsUsingFilter(bson.M{"component_name": name}, 1, 1, config.DefaultConfig.DbQueryTimeout)
}

func (c *ComponentStore) GetVulnerableSbomComponentsFilter(componentNames, componentVersions, sbomIds, compTypes, compNames, purls, versions []string, page, limit int) bson.M {
	conditions := map[string][]string{
		"component_name":    componentNames,
		"component_version": componentVersions,
		"purl":              purls,
		"sbom_id":           sbomIds,
		"type":              compTypes,
		"name":              compNames,
		"version":           versions,
	}
	filter := utils.BuildDynamicContainsFilter(conditions)

	// Add to the query: {vulns: { $exists: true, $ne: []}}
	filter["vulns"] = bson.M{
		"$exists": true,
		"$ne":     []interface{}{}, // Ensure the array is not empty
	}

	return filter
}

func (c *ComponentStore) GetVulnerableComponents(componentNames, componentVersions, sbomIds, compTypes, compNames, purls, versions []string, page, limit, duration int) (components []types.Component, total int64, err error) {
	filter := c.GetVulnerableSbomComponentsFilter(componentNames, componentVersions, sbomIds, compTypes, compNames, purls, versions, page, limit)

	components, err = c.GetComponentsUsingFilter(filter, page, limit, duration)
	if err != nil {
		log.Error().Err(err).Msg("failed to get components")
		return components, total, err
	}

	total, err = c.GetComponentTotalCount(filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get total components")
		return components, total, err
	}

	return components, total, err
}

func (c *ComponentStore) DeleteByIds(idParams []string, duration int) (int64, error) {
	// Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range idParams {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Error().Err(err).Msgf("Invalid ObjectID %s: %v", id, err)
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	// Define the filter to match any of the ObjectIDs
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	result, err := c.collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete documents: %v", err)
		return -1, err
	}

	return result.DeletedCount, nil
}

func (c *ComponentStore) DeleteById(idParam string, duration int) (int64, error) {
	return c.DeleteByIds([]string{idParam}, duration)
}
