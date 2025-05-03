// cloud/backend/internal/ipe/repo/person/presenter.go
package person

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

// FindPresenterByID retrieves a presenter by ID
func (impl personImpl) FindPresenterByID(ctx context.Context, id primitive.ObjectID) (*dom_person.Presenter, error) {
	filter := bson.M{"_id": id}

	var result dom_person.Presenter
	err := impl.PresenterCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get presenter by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindAllPresenters retrieves all presenters
func (impl personImpl) FindAllPresenters(ctx context.Context) ([]*dom_person.Presenter, error) {
	cursor, err := impl.PresenterCollection.Find(ctx, bson.M{})
	if err != nil {
		impl.Logger.Error("database find all presenters error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var presenters []*dom_person.Presenter
	if err = cursor.All(ctx, &presenters); err != nil {
		impl.Logger.Error("database decode presenters error", zap.Any("error", err))
		return nil, err
	}

	return presenters, nil
}

// SavePresenter creates a new presenter
func (impl personImpl) SavePresenter(ctx context.Context, presenter *dom_person.Presenter) (primitive.ObjectID, error) {
	if presenter.ID == primitive.NilObjectID {
		presenter.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert presenter without ID, created ID now", zap.Any("id", presenter.ID))
	}

	_, err := impl.PresenterCollection.InsertOne(ctx, presenter)
	if err != nil {
		impl.Logger.Error("database save presenter error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return presenter.ID, nil
}

// UpdatePresenter updates an existing presenter
func (impl personImpl) UpdatePresenter(ctx context.Context, presenter *dom_person.Presenter) error {
	filter := bson.M{"_id": presenter.ID}
	update := bson.M{"$set": presenter}

	_, err := impl.PresenterCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update presenter error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeletePresenter removes a presenter by ID
func (impl personImpl) DeletePresenter(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.PresenterCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete presenter error", zap.Any("error", err))
		return err
	}

	return nil
}
