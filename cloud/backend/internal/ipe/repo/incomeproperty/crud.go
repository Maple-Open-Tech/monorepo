// cloud/backend/internal/ipe/repo/incomeproperty/crud.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
)

// FindByID retrieves a property by its ID
func (impl incomePropertyImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*dom_property.IncomeProperty, error) {
	filter := bson.M{"_id": id}

	var result dom_property.IncomeProperty
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents
			return nil, nil
		}
		impl.Logger.Error("database get property by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindAll retrieves all properties
func (impl incomePropertyImpl) FindAll(ctx context.Context) ([]*dom_property.IncomeProperty, error) {
	cursor, err := impl.Collection.Find(ctx, bson.M{})
	if err != nil {
		impl.Logger.Error("database find all properties error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []*dom_property.IncomeProperty
	if err = cursor.All(ctx, &properties); err != nil {
		impl.Logger.Error("database decode properties error", zap.Any("error", err))
		return nil, err
	}

	return properties, nil
}

// Save creates a new property
func (impl incomePropertyImpl) Save(ctx context.Context, property *dom_property.IncomeProperty) (primitive.ObjectID, error) {
	if property.ID == primitive.NilObjectID {
		property.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert property without ID, created ID now", zap.Any("id", property.ID))
	}

	_, err := impl.Collection.InsertOne(ctx, property)
	if err != nil {
		impl.Logger.Error("database failed property creation", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return property.ID, nil
}

// Update updates an existing property
func (impl incomePropertyImpl) Update(ctx context.Context, property *dom_property.IncomeProperty) error {
	filter := bson.M{"_id": property.ID}
	update := bson.M{"$set": property}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update property error", zap.Any("error", err))
		return err
	}

	return nil
}

// Delete removes a property by ID
func (impl incomePropertyImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.Collection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete property error", zap.Any("error", err))
		return err
	}

	return nil
}

// FindByAddress finds properties by address
func (impl incomePropertyImpl) FindByAddress(ctx context.Context, address string) ([]dom_property.IncomeProperty, error) {
	filter := bson.M{"address": bson.M{"$regex": address, "$options": "i"}}

	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		impl.Logger.Error("database find properties by address error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []dom_property.IncomeProperty
	if err = cursor.All(ctx, &properties); err != nil {
		impl.Logger.Error("database decode properties error", zap.Any("error", err))
		return nil, err
	}

	return properties, nil
}

// FindByCity finds properties by city
func (impl incomePropertyImpl) FindByCity(ctx context.Context, city string) ([]dom_property.IncomeProperty, error) {
	filter := bson.M{"city": bson.M{"$regex": city, "$options": "i"}}

	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		impl.Logger.Error("database find properties by city error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []dom_property.IncomeProperty
	if err = cursor.All(ctx, &properties); err != nil {
		impl.Logger.Error("database decode properties error", zap.Any("error", err))
		return nil, err
	}

	return properties, nil
}
