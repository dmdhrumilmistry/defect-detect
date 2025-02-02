package sbom

import (
	"context"
	"fmt"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COMPONENT_SBOM_COLLECTION = "component_sbom"

type ComponentSbomStore struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewComponentSbomStore(db *mongo.Database) *ComponentSbomStore {
	collection := db.Collection(COMPONENT_SBOM_COLLECTION)
	// TODO: create index if not exists

	return &ComponentSbomStore{
		db:         db,
		collection: collection,
	}
}

func (c *ComponentSbomStore) AddComponentSbom(sbom cyclonedx.BOM) (string, error) {
	result, err := c.collection.InsertOne(context.TODO(), sbom)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert sbom")
		return "", err
	}

	return (result.InsertedID).(primitive.ObjectID).Hex(), nil
}

func (c *ComponentSbomStore) ValidateIds(ids []string) error {
	// Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Error().Err(err).Msgf("Invalid ObjectID %s: %v", id, err)
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	filter := bson.M{"_id": bson.M{
		"$in": objectIDs,
	}}

	count, err := c.GetTotalCount(filter)
	if err != nil {
		return err
	}

	expectedLen := len(ids)
	totalCount := int(count)
	if expectedLen != totalCount {
		return fmt.Errorf("expected sbom count to be %d got %d", expectedLen, totalCount)
	}

	return nil
}

func (c *ComponentSbomStore) GetTotalCount(filter interface{}) (int64, error) {
	// Get total count of documents
	total, err := c.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Handler for getting paginated items
func (c *ComponentSbomStore) GetPaginatedSboms(page, limit, duration int) ([]types.Sbom, error) {
	var sboms []types.Sbom

	// Calculate skip
	skip := (page - 1) * limit

	// MongoDB query options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	cursor, err := c.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return sboms, err
	}
	defer cursor.Close(ctx)

	// Parse results
	if err := cursor.All(ctx, &sboms); err != nil {
		return sboms, err
	}

	return sboms, nil
}

// Handler for getting paginated items
func (c *ComponentSbomStore) GetSbomById(idParam string, duration int) (types.Sbom, error) {
	var sbom types.Sbom

	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return sbom, err
	}

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	err = c.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&sbom)
	if err != nil {
		return sbom, err
	}

	return sbom, nil
}

func (c *ComponentSbomStore) GetSbomByName(name string, duration int) ([]types.Sbom, error) {
	var sboms []types.Sbom

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	cursor, err := c.collection.Find(ctx, bson.M{"metadata.component.name": name})
	if err != nil {
		return sboms, err
	}
	defer cursor.Close(ctx)

	// Parse results
	if err := cursor.All(ctx, &sboms); err != nil {
		return sboms, err
	}

	return sboms, nil
}

func (c *ComponentSbomStore) DeleteByIds(idParams []string, duration int) (int64, error) {
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

func (c *ComponentSbomStore) DeleteById(idParam string, duration int) (int64, error) {
	return c.DeleteByIds([]string{idParam}, duration)
}
