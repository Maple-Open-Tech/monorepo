// cloud/backend/internal/ipe/repo/evaluation/photo.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"

	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
)

// AddPropertyPhoto adds a property photo
func (impl evaluationImpl) AddPropertyPhoto(ctx context.Context, evaluationID primitive.ObjectID, photo *dom_evaluation.PropertyPhoto) error {
	if photo.ID == primitive.NilObjectID {
		photo.ID = primitive.NewObjectID()
	}
	photo.EvaluationID = evaluationID

	_, err := impl.PropertyPhotoCollection.InsertOne(ctx, photo)
	if err != nil {
		impl.Logger.Error("database add property photo error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdatePropertyPhoto updates a property photo
func (impl evaluationImpl) UpdatePropertyPhoto(ctx context.Context, photo *dom_evaluation.PropertyPhoto) error {
	filter := bson.M{"_id": photo.ID}
	update := bson.M{"$set": photo}

	_, err := impl.PropertyPhotoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update property photo error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeletePropertyPhoto removes a property photo
func (impl evaluationImpl) DeletePropertyPhoto(ctx context.Context, photoID primitive.ObjectID) error {
	filter := bson.M{"_id": photoID}

	_, err := impl.PropertyPhotoCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete property photo error", zap.Any("error", err))
		return err
	}

	return nil
}

// FindPhotosByEvaluationID retrieves photos by evaluation ID
func (impl evaluationImpl) FindPhotosByEvaluationID(ctx context.Context, evaluationID primitive.ObjectID) ([]*dom_evaluation.PropertyPhoto, error) {
	filter := bson.M{"evaluationId": evaluationID}

	cursor, err := impl.PropertyPhotoCollection.Find(ctx, filter)
	if err != nil {
		impl.Logger.Error("database find photos by evaluation id error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var photos []*dom_evaluation.PropertyPhoto
	if err = cursor.All(ctx, &photos); err != nil {
		impl.Logger.Error("database decode photos error", zap.Any("error", err))
		return nil, err
	}

	return photos, nil
}
