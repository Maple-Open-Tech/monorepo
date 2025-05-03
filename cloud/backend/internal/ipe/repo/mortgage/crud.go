// cloud/backend/internal/ipe/repo/mortgage/crud.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
)

// FindByID retrieves a mortgage by ID
func (impl mortgageImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*dom_mortgage.Mortgage, error) {
	filter := bson.M{"_id": id}

	var result dom_mortgage.Mortgage
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get mortgage by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindByFinancialAnalysisID retrieves a mortgage by financial analysis ID
func (impl mortgageImpl) FindByFinancialAnalysisID(ctx context.Context, analysisID primitive.ObjectID) (*dom_mortgage.Mortgage, error) {
	filter := bson.M{"financialAnalysisId": analysisID}

	var result dom_mortgage.Mortgage
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get mortgage by financial analysis id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// Save creates a new mortgage
func (impl mortgageImpl) Save(ctx context.Context, mortgage *dom_mortgage.Mortgage) (primitive.ObjectID, error) {
	if mortgage.ID == primitive.NilObjectID {
		mortgage.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert mortgage without ID, created ID now", zap.Any("id", mortgage.ID))
	}

	_, err := impl.Collection.InsertOne(ctx, mortgage)
	if err != nil {
		impl.Logger.Error("database save mortgage error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return mortgage.ID, nil
}

// Update updates an existing mortgage
func (impl mortgageImpl) Update(ctx context.Context, mortgage *dom_mortgage.Mortgage) error {
	filter := bson.M{"_id": mortgage.ID}
	update := bson.M{"$set": mortgage}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update mortgage error", zap.Any("error", err))
		return err
	}

	return nil
}

// Delete removes a mortgage by ID
func (impl mortgageImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.Collection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete mortgage error", zap.Any("error", err))
		return err
	}

	return nil
}
