// cloud/backend/internal/ipe/repo/financialanalysis/income.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"

	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
)

// AddRentalIncome adds rental income to financial analysis
func (impl financialAnalysisImpl) AddRentalIncome(ctx context.Context, analysisID primitive.ObjectID, income *dom_financial.RentalIncome) error {
	if income.ID == primitive.NilObjectID {
		income.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"rentalIncomes": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add rental income error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateRentalIncome updates a rental income
func (impl financialAnalysisImpl) UpdateRentalIncome(ctx context.Context, income *dom_financial.RentalIncome) error {
	filter := bson.M{"rentalIncomes._id": income.ID}
	update := bson.M{"$set": bson.M{"rentalIncomes.$": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update rental income error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteRentalIncome removes a rental income
func (impl financialAnalysisImpl) DeleteRentalIncome(ctx context.Context, incomeID primitive.ObjectID) error {
	filter := bson.M{"rentalIncomes._id": incomeID}
	update := bson.M{"$pull": bson.M{"rentalIncomes": bson.M{"_id": incomeID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete rental income error", zap.Any("error", err))
		return err
	}

	return nil
}

// AddCommercialIncome adds commercial income to financial analysis
func (impl financialAnalysisImpl) AddCommercialIncome(ctx context.Context, analysisID primitive.ObjectID, income *dom_financial.CommercialIncome) error {
	if income.ID == primitive.NilObjectID {
		income.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"commercialIncomes": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add commercial income error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateCommercialIncome updates a commercial income
func (impl financialAnalysisImpl) UpdateCommercialIncome(ctx context.Context, income *dom_financial.CommercialIncome) error {
	filter := bson.M{"commercialIncomes._id": income.ID}
	update := bson.M{"$set": bson.M{"commercialIncomes.$": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update commercial income error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteCommercialIncome removes a commercial income
func (impl financialAnalysisImpl) DeleteCommercialIncome(ctx context.Context, incomeID primitive.ObjectID) error {
	filter := bson.M{"commercialIncomes._id": incomeID}
	update := bson.M{"$pull": bson.M{"commercialIncomes": bson.M{"_id": incomeID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete commercial income error", zap.Any("error", err))
		return err
	}

	return nil
}

// AddFacilityIncome adds facility income to financial analysis
func (impl financialAnalysisImpl) AddFacilityIncome(ctx context.Context, analysisID primitive.ObjectID, income *dom_financial.FacilityIncome) error {
	if income.ID == primitive.NilObjectID {
		income.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"facilityIncomes": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add facility income error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateFacilityIncome updates a facility income
func (impl financialAnalysisImpl) UpdateFacilityIncome(ctx context.Context, income *dom_financial.FacilityIncome) error {
	filter := bson.M{"facilityIncomes._id": income.ID}
	update := bson.M{"$set": bson.M{"facilityIncomes.$": income}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update facility income error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteFacilityIncome removes a facility income
func (impl financialAnalysisImpl) DeleteFacilityIncome(ctx context.Context, incomeID primitive.ObjectID) error {
	filter := bson.M{"facilityIncomes._id": incomeID}
	update := bson.M{"$pull": bson.M{"facilityIncomes": bson.M{"_id": incomeID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete facility income error", zap.Any("error", err))
		return err
	}

	return nil
}
