package component

import (
	"context"
	"fmt"
	"time"

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
}

func NewComponentStore(db *mongo.Database) *ComponentStore {
	collection := db.Collection(COMPONENT_COLLECTION)
	// TODO: create index if not exists

	return &ComponentStore{
		db:         db,
		collection: collection,
	}
}

func (c *ComponentStore) AddComponentUsingSbom(sbom types.Sbom) ([]string, error) {
	componentName := sbom.Metadata.Component.Name
	componentVersion := sbom.Metadata.Component.Version
	insertedIds := []string{}

	var components []interface{}
	for _, component := range *sbom.Components {
		log.Print(component)
		log.Print(component.Licenses)
		fmt.Print("=================")

		// fetch licenses slice from sbom
		var licences []string
		if component.Licenses != nil {
			for _, license := range *component.Licenses {
				if license.License != nil && license.License.ID != "" {
					log.Print(license.License.ID)
					licences = append(licences, license.License.ID)
				}
			}
		}

		// create components slice
		components = append(components, types.Component{
			Name:             component.Name,
			Version:          component.Version,
			PackageUrl:       component.PackageURL,
			Licenses:         licences,
			Type:             string(component.Type),
			ComponentName:    componentName,
			ComponentVersion: componentVersion,
		})
	}

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
