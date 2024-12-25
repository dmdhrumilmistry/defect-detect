package component

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COMPONENT_COLLECTION = "component"

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

func (c *ComponentStore) processComponents(sbom types.Sbom, componentName, componentVersion string, workers int) []interface{} {
	var components []interface{}
	type vulnResult struct {
		Component types.Component
		Err       error
	}

	// Channels for work distribution and results collection
	workCh := make(chan *cyclonedx.Component)
	resultCh := make(chan vulnResult)

	// Worker function
	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for component := range workCh {
			var licences []string
			if component.Licenses != nil {
				for _, license := range *component.Licenses {
					if license.License != nil && license.License.ID != "" {
						licences = append(licences, license.License.ID)
					}
				}
			}

			var vulns []types.Vuln
			var err error
			if component.PackageURL != "" {
				vulns, err = c.Analyzer.GetVulns(component.PackageURL)
				if err != nil {
					log.Error().Err(err).Msgf("failed to analyze vulns for %s", component.PackageURL)
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
				},
				Err: err,
			}
		}
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(&wg)
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

	components := c.processComponents(sbom, componentName, componentVersion, 20)

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

func (c *ComponentStore) GetComponentTotalCount() (int64, error) {
	// Get total count of documents
	total, err := c.collection.CountDocuments(context.TODO(), bson.M{})
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

	return c.GetComponentsUsingFilter(bson.M{"_id": objID}, 1, 1, 5)
}

func (c *ComponentStore) GetComponentByName(name string, duration int) ([]types.Component, error) {
	return c.GetComponentsUsingFilter(bson.M{"component_name": name}, 1, 1, 5)
}
