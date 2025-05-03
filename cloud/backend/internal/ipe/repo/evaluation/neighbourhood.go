// cloud/backend/internal/ipe/repo/evaluation/neighbourhood.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
)

// FindNeighbourhoodByID retrieves a neighbourhood by ID
func (impl evaluationImpl) FindNeighbourhoodByID(ctx context.Context, id primitive.ObjectID) (*dom_evaluation.Neighbourhood, error) {
	filter := bson.M{"_id": id}

	var result dom_evaluation.Neighbourhood
	err := impl.NeighbourhoodCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get neighbourhood by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// SaveNeighbourhood creates a new neighbourhood
func (impl evaluationImpl) SaveNeighbourhood(ctx context.Context, neighbourhood *dom_evaluation.Neighbourhood) (primitive.ObjectID, error) {
	if neighbourhood.ID == primitive.NilObjectID {
		neighbourhood.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert neighbourhood without ID, created ID now", zap.Any("id", neighbourhood.ID))
	}

	_, err := impl.NeighbourhoodCollection.InsertOne(ctx, neighbourhood)
	if err != nil {
		impl.Logger.Error("database save neighbourhood error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return neighbourhood.ID, nil
}

// UpdateNeighbourhood updates an existing neighbourhood
func (impl evaluationImpl) UpdateNeighbourhood(ctx context.Context, neighbourhood *dom_evaluation.Neighbourhood) error {
	filter := bson.M{"_id": neighbourhood.ID}
	update := bson.M{"$set": neighbourhood}

	_, err := impl.NeighbourhoodCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update neighbourhood error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteNeighbourhood removes a neighbourhood by ID
func (impl evaluationImpl) DeleteNeighbourhood(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.NeighbourhoodCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete neighbourhood error", zap.Any("error", err))
		return err
	}

	return nil
}
