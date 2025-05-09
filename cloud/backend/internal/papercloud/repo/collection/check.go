// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/collection/check.go
package collection

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (impl collectionStorerImpl) CheckIfExistsByID(id string) (bool, error) {
	ctx := context.Background()
	filter := bson.M{"id": id}

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by ID error", zap.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}

func (impl collectionStorerImpl) CheckIfUserHasAccess(collectionID string, userID string) (bool, error) {
	ctx := context.Background()

	// User has access if they own the collection or it's shared with them
	filter := bson.M{
		"id": collectionID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"shared_with.user_id": userID},
		},
	}

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if user has access error", zap.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
