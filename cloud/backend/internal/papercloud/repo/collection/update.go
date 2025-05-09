// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/collection/update.go
package collection

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"

	dom_collection "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/collection"
)

func (impl collectionStorerImpl) Update(collection *dom_collection.Collection) error {
	ctx := context.Background()
	filter := bson.M{"id": collection.ID}

	update := bson.M{
		"$set": collection,
	}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update collection error",
			zap.Any("error", err),
			zap.String("id", collection.ID))
		return err
	}

	return nil
}

func (impl collectionStorerImpl) AddShare(collectionID string, share *dom_collection.Share) error {
	ctx := context.Background()
	filter := bson.M{"id": collectionID}

	// Add share to the shared_with array
	update := bson.M{
		"$push": bson.M{
			"shared_with": share,
		},
	}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add share error",
			zap.Any("error", err),
			zap.String("collection_id", collectionID),
			zap.String("user_id", share.UserID))
		return err
	}

	return nil
}

func (impl collectionStorerImpl) RemoveShare(collectionID string, userID string) error {
	ctx := context.Background()
	filter := bson.M{"id": collectionID}

	// Remove the share from the shared_with array
	update := bson.M{
		"$pull": bson.M{
			"shared_with": bson.M{
				"user_id": userID,
			},
		},
	}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database remove share error",
			zap.Any("error", err),
			zap.String("collection_id", collectionID),
			zap.String("user_id", userID))
		return err
	}

	return nil
}
