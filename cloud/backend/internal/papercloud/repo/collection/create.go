// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/collection/create.go
package collection

import (
	"context"
	"errors"

	"go.uber.org/zap"

	dom_collection "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/collection"
)

func (impl collectionStorerImpl) Create(collection *dom_collection.Collection) error {
	ctx := context.Background()

	// Validate collection ID
	if collection.ID == "" {
		impl.Logger.Error("collection ID is required but not provided")
		return errors.New("collection ID is required")
	}

	// Insert collection document
	_, err := impl.Collection.InsertOne(ctx, collection)

	if err != nil {
		impl.Logger.Error("database failed create collection error",
			zap.Any("error", err),
			zap.String("id", collection.ID))
		return err
	}

	return nil
}
