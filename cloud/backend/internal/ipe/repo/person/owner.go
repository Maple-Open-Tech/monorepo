// cloud/backend/internal/ipe/repo/person/owner.go
package person

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

// FindOwnerByID retrieves an owner by ID
func (impl personImpl) FindOwnerByID(ctx context.Context, id primitive.ObjectID) (*dom_person.Owner, error) {
	filter := bson.M{"_id": id}

	var result dom_person.Owner
	err := impl.OwnerCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get owner by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindAllOwners retrieves all owners
func (impl personImpl) FindAllOwners(ctx context.Context) ([]*dom_person.Owner, error) {
	cursor, err := impl.OwnerCollection.Find(ctx, bson.M{})
	if err != nil {
		impl.Logger.Error("database find all owners error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var owners []*dom_person.Owner
	if err = cursor.All(ctx, &owners); err != nil {
		impl.Logger.Error("database decode owners error", zap.Any("error", err))
		return nil, err
	}

	return owners, nil
}

// SaveOwner creates a new owner
func (impl personImpl) SaveOwner(ctx context.Context, owner *dom_person.Owner) (primitive.ObjectID, error) {
	if owner.ID == primitive.NilObjectID {
		owner.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert owner without ID, created ID now", zap.Any("id", owner.ID))
	}

	_, err := impl.OwnerCollection.InsertOne(ctx, owner)
	if err != nil {
		impl.Logger.Error("database save owner error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return owner.ID, nil
}

// UpdateOwner updates an existing owner
func (impl personImpl) UpdateOwner(ctx context.Context, owner *dom_person.Owner) error {
	filter := bson.M{"_id": owner.ID}
	update := bson.M{"$set": owner}

	_, err := impl.OwnerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update owner error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteOwner removes an owner by ID
func (impl personImpl) DeleteOwner(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.OwnerCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete owner error", zap.Any("error", err))
		return err
	}

	return nil
}
