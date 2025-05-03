// cloud/backend/internal/ipe/repo/financialanalysis/other.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"

	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
)

// AddAnnualProjection adds an annual projection
func (impl financialAnalysisImpl) AddAnnualProjection(ctx context.Context, analysisID primitive.ObjectID, projection *dom_financial.AnnualProjection) error {
	if projection.ID == primitive.NilObjectID {
		projection.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"annualProjections": projection}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add annual projection error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateAnnualProjection updates an annual projection
func (impl financialAnalysisImpl) UpdateAnnualProjection(ctx context.Context, projection *dom_financial.AnnualProjection) error {
	filter := bson.M{"annualProjections._id": projection.ID}
	update := bson.M{"$set": bson.M{"annualProjections.$": projection}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update annual projection error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteAnnualProjection removes an annual projection
func (impl financialAnalysisImpl) DeleteAnnualProjection(ctx context.Context, projectionID primitive.ObjectID) error {
	filter := bson.M{"annualProjections._id": projectionID}
	update := bson.M{"$pull": bson.M{"annualProjections": bson.M{"_id": projectionID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete annual projection error", zap.Any("error", err))
		return err
	}

	return nil
}

// AddPurchaseFee adds a purchase fee
func (impl financialAnalysisImpl) AddPurchaseFee(ctx context.Context, analysisID primitive.ObjectID, fee *dom_financial.PurchaseFee) error {
	if fee.ID == primitive.NilObjectID {
		fee.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"purchaseFees": fee}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add purchase fee error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdatePurchaseFee updates a purchase fee
func (impl financialAnalysisImpl) UpdatePurchaseFee(ctx context.Context, fee *dom_financial.PurchaseFee) error {
	filter := bson.M{"purchaseFees._id": fee.ID}
	update := bson.M{"$set": bson.M{"purchaseFees.$": fee}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update purchase fee error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeletePurchaseFee removes a purchase fee
func (impl financialAnalysisImpl) DeletePurchaseFee(ctx context.Context, feeID primitive.ObjectID) error {
	filter := bson.M{"purchaseFees._id": feeID}
	update := bson.M{"$pull": bson.M{"purchaseFees": bson.M{"_id": feeID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete purchase fee error", zap.Any("error", err))
		return err
	}

	return nil
}

// AddCapitalImprovement adds a capital improvement
func (impl financialAnalysisImpl) AddCapitalImprovement(ctx context.Context, analysisID primitive.ObjectID, improvement *dom_financial.CapitalImprovement) error {
	if improvement.ID == primitive.NilObjectID {
		improvement.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"capitalImprovements": improvement}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add capital improvement error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateCapitalImprovement updates a capital improvement
func (impl financialAnalysisImpl) UpdateCapitalImprovement(ctx context.Context, improvement *dom_financial.CapitalImprovement) error {
	filter := bson.M{"capitalImprovements._id": improvement.ID}
	update := bson.M{"$set": bson.M{"capitalImprovements.$": improvement}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update capital improvement error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteCapitalImprovement removes a capital improvement
func (impl financialAnalysisImpl) DeleteCapitalImprovement(ctx context.Context, improvementID primitive.ObjectID) error {
	filter := bson.M{"capitalImprovements._id": improvementID}
	update := bson.M{"$pull": bson.M{"capitalImprovements": bson.M{"_id": improvementID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete capital improvement error", zap.Any("error", err))
		return err
	}

	return nil
}
