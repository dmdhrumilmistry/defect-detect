package project

import (
	"context"
	"time"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PROJECT_COLLECTION = "project"

type ProjectStore struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProjectStore(db *mongo.Database) *ProjectStore {
	collection := db.Collection(PROJECT_COLLECTION)
	return &ProjectStore{
		db:         db,
		collection: collection,
	}
}

func (p *ProjectStore) AddProject(project types.Project) (string, error) {
	result, err := p.collection.InsertOne(context.TODO(), project)
	if err != nil {
		log.Error().Err(err).Msg("failed to ins")
		return "", err
	}

	return (result.InsertedID).(primitive.ObjectID).Hex(), nil
}

func (p *ProjectStore) GetUsingFilter(filter interface{}, page, limit, duration int) ([]types.Project, error) {
	var projects []types.Project

	// Calculate skip
	skip := (page - 1) * limit

	// MongoDB query options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	cursor, err := p.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return projects, err
	}
	defer cursor.Close(ctx)

	// Parse results
	if err := cursor.All(ctx, &projects); err != nil {
		return projects, err
	}

	return projects, nil
}

// Handler for getting paginated items
func (p *ProjectStore) GetProjectById(idParam string, duration int) ([]types.Project, error) {
	// Convert the string ID to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return []types.Project{}, err
	}

	return p.GetUsingFilter(bson.M{"_id": objID}, 1, 1, config.DefaultConfig.DbQueryTimeout)
}

func (p *ProjectStore) GetByName(name string, duration int) ([]types.Project, error) {
	return p.GetUsingFilter(bson.M{"name": name}, 1, 1, config.DefaultConfig.DbQueryTimeout)
}

func (p *ProjectStore) DeleteByIds(idParams []string, duration int) (int64, error) {
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

	result, err := p.collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete documents: %v", err)
		return -1, err
	}

	return result.DeletedCount, nil
}

func (p *ProjectStore) DeleteById(idParam string, duration int) (int64, error) {
	return p.DeleteByIds([]string{idParam}, duration)
}

func (p *ProjectStore) GetTotalCount(filter interface{}) (int64, error) {
	// Get total count of documents
	total, err := p.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return total, nil
}
