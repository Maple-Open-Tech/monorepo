// cloud/backend/internal/ipe/repo/mortgage/intervals.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"

	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
)

// AddMortgageInterval adds a mortgage interval to a mortgage
func (impl mortgageImpl) AddMortgageInterval(ctx context.Context, mortgageID primitive.ObjectID, interval *dom_mortgage.MortgageInterval) error {
	if interval.ID == primitive.NilObjectID {
		interval.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": mortgageID}
	update := bson.M{"$push": bson.M{"mortgagePaymentSchedule": interval}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add mortgage interval error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateMortgageInterval updates a mortgage interval
func (impl mortgageImpl) UpdateMortgageInterval(ctx context.Context, interval *dom_mortgage.MortgageInterval) error {
	filter := bson.M{"mortgagePaymentSchedule._id": interval.ID}
	update := bson.M{"$set": bson.M{"mortgagePaymentSchedule.$": interval}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update mortgage interval error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteMortgageInterval removes a mortgage interval
func (impl mortgageImpl) DeleteMortgageInterval(ctx context.Context, intervalID primitive.ObjectID) error {
	filter := bson.M{"mortgagePaymentSchedule._id": intervalID}
	update := bson.M{"$pull": bson.M{"mortgagePaymentSchedule": bson.M{"_id": intervalID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete mortgage interval error", zap.Any("error", err))
		return err
	}

	return nil
}

// ClearMortgageIntervals removes all mortgage intervals from a mortgage
func (impl mortgageImpl) ClearMortgageIntervals(ctx context.Context, mortgageID primitive.ObjectID) error {
	filter := bson.M{"_id": mortgageID}
	update := bson.M{"$set": bson.M{"mortgagePaymentSchedule": []dom_mortgage.MortgageInterval{}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database clear mortgage intervals error", zap.Any("error", err))
		return err
	}

	return nil
}
