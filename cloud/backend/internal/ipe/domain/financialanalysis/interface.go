package financialanalysis

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FinancialRepository defines methods for financial analysis storage operations
type FinancialRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*FinancialAnalysis, error)
	FindByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*FinancialAnalysis, error)
	Save(ctx context.Context, analysis *FinancialAnalysis) (primitive.ObjectID, error)
	Update(ctx context.Context, analysis *FinancialAnalysis) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Income methods
	AddRentalIncome(ctx context.Context, analysisID primitive.ObjectID, income *RentalIncome) error
	UpdateRentalIncome(ctx context.Context, income *RentalIncome) error
	DeleteRentalIncome(ctx context.Context, incomeID primitive.ObjectID) error

	AddCommercialIncome(ctx context.Context, analysisID primitive.ObjectID, income *CommercialIncome) error
	UpdateCommercialIncome(ctx context.Context, income *CommercialIncome) error
	DeleteCommercialIncome(ctx context.Context, incomeID primitive.ObjectID) error

	AddFacilityIncome(ctx context.Context, analysisID primitive.ObjectID, income *FacilityIncome) error
	UpdateFacilityIncome(ctx context.Context, income *FacilityIncome) error
	DeleteFacilityIncome(ctx context.Context, incomeID primitive.ObjectID) error

	// Expense methods
	AddExpense(ctx context.Context, analysisID primitive.ObjectID, expense *Expense) error
	UpdateExpense(ctx context.Context, expense *Expense) error
	DeleteExpense(ctx context.Context, expenseID primitive.ObjectID) error

	// Other methods
	AddAnnualProjection(ctx context.Context, analysisID primitive.ObjectID, projection *AnnualProjection) error
	UpdateAnnualProjection(ctx context.Context, projection *AnnualProjection) error
	DeleteAnnualProjection(ctx context.Context, projectionID primitive.ObjectID) error

	AddPurchaseFee(ctx context.Context, analysisID primitive.ObjectID, fee *PurchaseFee) error
	UpdatePurchaseFee(ctx context.Context, fee *PurchaseFee) error
	DeletePurchaseFee(ctx context.Context, feeID primitive.ObjectID) error

	AddCapitalImprovement(ctx context.Context, analysisID primitive.ObjectID, improvement *CapitalImprovement) error
	UpdateCapitalImprovement(ctx context.Context, improvement *CapitalImprovement) error
	DeleteCapitalImprovement(ctx context.Context, improvementID primitive.ObjectID) error
}
