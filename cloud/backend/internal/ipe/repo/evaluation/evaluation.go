// cloud/backend/internal/ipe/repo/evaluation/evaluation.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
)

// FindByID retrieves an evaluation by ID
func (impl evaluationImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*dom_evaluation.Evaluation, error) {
	filter := bson.M{"_id": id}

	var result dom_evaluation.Evaluation
	err := impl.EvaluationCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get evaluation by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindByPropertyID retrieves an evaluation by property ID
func (impl evaluationImpl) FindByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*dom_evaluation.Evaluation, error) {
	filter := bson.M{"propertyId": propertyID}

	var result dom_evaluation.Evaluation
	err := impl.EvaluationCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get evaluation by property id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// Save creates a new evaluation
func (impl evaluationImpl) Save(ctx context.Context, evaluation *dom_evaluation.Evaluation) (primitive.ObjectID, error) {
	if evaluation.ID == primitive.NilObjectID {
		evaluation.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert evaluation without ID, created ID now", zap.Any("id", evaluation.ID))
	}

	_, err := impl.EvaluationCollection.InsertOne(ctx, evaluation)
	if err != nil {
		impl.Logger.Error("database save evaluation error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return evaluation.ID, nil
}

// Update updates an existing evaluation
func (impl evaluationImpl) Update(ctx context.Context, evaluation *dom_evaluation.Evaluation) error {
	filter := bson.M{"_id": evaluation.ID}
	update := bson.M{"$set": evaluation}

	_, err := impl.EvaluationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update evaluation error", zap.Any("error", err))
		return err
	}

	return nil
}

// Delete removes an evaluation by ID
func (impl evaluationImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.EvaluationCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete evaluation error", zap.Any("error", err))
		return err
	}

	return nil
}
