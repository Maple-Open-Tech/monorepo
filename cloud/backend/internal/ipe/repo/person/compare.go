// cloud/backend/internal/ipe/repo/person/compare.go
package person

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

// FindCompareByID retrieves a compare by ID
func (impl personImpl) FindCompareByID(ctx context.Context, id primitive.ObjectID) (*dom_person.Compare, error) {
	filter := bson.M{"_id": id}

	var result dom_person.Compare
	err := impl.CompareCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get compare by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// SaveCompare creates a new compare
func (impl personImpl) SaveCompare(ctx context.Context, compare *dom_person.Compare) (primitive.ObjectID, error) {
	if compare.ID == primitive.NilObjectID {
		compare.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert compare without ID, created ID now", zap.Any("id", compare.ID))
	}

	_, err := impl.CompareCollection.InsertOne(ctx, compare)
	if err != nil {
		impl.Logger.Error("database save compare error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return compare.ID, nil
}

// UpdateCompare updates an existing compare
func (impl personImpl) UpdateCompare(ctx context.Context, compare *dom_person.Compare) error {
	filter := bson.M{"_id": compare.ID}
	update := bson.M{"$set": compare}

	_, err := impl.CompareCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update compare error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteCompare removes a compare by ID
func (impl personImpl) DeleteCompare(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.CompareCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete compare error", zap.Any("error", err))
		return err
	}

	return nil
}
