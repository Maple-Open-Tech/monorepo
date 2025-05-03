// cloud/backend/internal/ipe/repo/financialanalysis/expense.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"

	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
)

// AddExpense adds an expense to financial analysis
func (impl financialAnalysisImpl) AddExpense(ctx context.Context, analysisID primitive.ObjectID, expense *dom_financial.Expense) error {
	if expense.ID == primitive.NilObjectID {
		expense.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": analysisID}
	update := bson.M{"$push": bson.M{"expenses": expense}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database add expense error", zap.Any("error", err))
		return err
	}

	return nil
}

// UpdateExpense updates an expense
func (impl financialAnalysisImpl) UpdateExpense(ctx context.Context, expense *dom_financial.Expense) error {
	filter := bson.M{"expenses._id": expense.ID}
	update := bson.M{"$set": bson.M{"expenses.$": expense}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update expense error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteExpense removes an expense
func (impl financialAnalysisImpl) DeleteExpense(ctx context.Context, expenseID primitive.ObjectID) error {
	filter := bson.M{"expenses._id": expenseID}
	update := bson.M{"$pull": bson.M{"expenses": bson.M{"_id": expenseID}}}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database delete expense error", zap.Any("error", err))
		return err
	}

	return nil
}
