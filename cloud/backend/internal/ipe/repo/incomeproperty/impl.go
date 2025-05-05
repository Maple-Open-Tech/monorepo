package incomeproperty

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"

	dom "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
)

type mongoRepository struct {
	db                 *mongo.Database
	logger             *zap.Logger
	propertyCollection *mongo.Collection
}

// NewMongoRepository creates a new MongoDB repository for the Income Property Evaluator
func NewMongoRepository(db *mongo.Database, logger *zap.Logger) dom.Repository {
	return &mongoRepository{
		db:                 db,
		logger:             logger,
		propertyCollection: db.Collection("income_properties"),
	}
}

// Find a property by ID
func (r *mongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*dom.IncomeProperty, error) {
	property := &dom.IncomeProperty{}
	err := r.propertyCollection.FindOne(ctx, bson.M{"_id": id}).Decode(property)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("property not found: %w", err)
		}
		return nil, fmt.Errorf("error finding property: %w", err)
	}
	return property, nil
}

// Find all properties
func (r *mongoRepository) FindAll(ctx context.Context) ([]*dom.IncomeProperty, error) {
	cursor, err := r.propertyCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding properties: %w", err)
	}
	defer cursor.Close(ctx)

	properties := []*dom.IncomeProperty{}
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, fmt.Errorf("error decoding properties: %w", err)
	}
	return properties, nil
}

// Create a property
func (r *mongoRepository) Create(ctx context.Context, property *dom.IncomeProperty) (primitive.ObjectID, error) {
	if property.ID.IsZero() {
		property.ID = primitive.NewObjectID()
	}

	result, err := r.propertyCollection.InsertOne(ctx, property)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error saving property: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// Update a property
func (r *mongoRepository) Update(ctx context.Context, property *dom.IncomeProperty) error {
	if property.ID.IsZero() {
		return errors.New("property ID is required for update")
	}

	result, err := r.propertyCollection.ReplaceOne(ctx, bson.M{"_id": property.ID}, property)
	if err != nil {
		return fmt.Errorf("error updating property: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("property not found with ID: %s", property.ID.Hex())
	}

	return nil
}

// Delete a property
func (r *mongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.propertyCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting property: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("property not found with ID: %s", id.Hex())
	}

	return nil
}

// Find properties by address
func (r *mongoRepository) FindByAddress(ctx context.Context, address string) ([]*dom.IncomeProperty, error) {
	cursor, err := r.propertyCollection.Find(ctx, bson.M{"address": bson.M{"$regex": address, "$options": "i"}})
	if err != nil {
		return nil, fmt.Errorf("error finding properties by address: %w", err)
	}
	defer cursor.Close(ctx)

	properties := []*dom.IncomeProperty{}
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, fmt.Errorf("error decoding properties: %w", err)
	}
	return properties, nil
}

// Find properties by city
func (r *mongoRepository) FindByCity(ctx context.Context, city string) ([]*dom.IncomeProperty, error) {
	cursor, err := r.propertyCollection.Find(ctx, bson.M{"city": bson.M{"$regex": city, "$options": "i"}})
	if err != nil {
		return nil, fmt.Errorf("error finding properties by city: %w", err)
	}
	defer cursor.Close(ctx)

	properties := []*dom.IncomeProperty{}
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, fmt.Errorf("error decoding properties: %w", err)
	}
	return properties, nil
}
