// cloud/backend/internal/ipe/repo/financialanalysis/crud.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
)

// FindByID retrieves a financial analysis by ID
func (impl financialAnalysisImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*dom_financial.FinancialAnalysis, error) {
	filter := bson.M{"_id": id}

	var result dom_financial.FinancialAnalysis
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get financial analysis by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindByPropertyID retrieves a financial analysis by property ID
func (impl financialAnalysisImpl) FindByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*dom_financial.FinancialAnalysis, error) {
	filter := bson.M{"propertyId": propertyID}

	var result dom_financial.FinancialAnalysis
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get financial analysis by property id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// Save creates a new financial analysis
func (impl financialAnalysisImpl) Save(ctx context.Context, analysis *dom_financial.FinancialAnalysis) (primitive.ObjectID, error) {
	if analysis.ID == primitive.NilObjectID {
		analysis.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert financial analysis without ID, created ID now", zap.Any("id", analysis.ID))
	}

	_, err := impl.Collection.InsertOne(ctx, analysis)
	if err != nil {
		impl.Logger.Error("database save financial analysis error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return analysis.ID, nil
}

// Update updates an existing financial analysis
func (impl financialAnalysisImpl) Update(ctx context.Context, analysis *dom_financial.FinancialAnalysis) error {
	filter := bson.M{"_id": analysis.ID}
	update := bson.M{"$set": analysis}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update financial analysis error", zap.Any("error", err))
		return err
	}

	return nil
}

// Delete removes a financial analysis by ID
func (impl financialAnalysisImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.Collection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete financial analysis error", zap.Any("error", err))
		return err
	}

	return nil
}
