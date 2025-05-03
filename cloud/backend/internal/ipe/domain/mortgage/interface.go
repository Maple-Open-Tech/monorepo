package mortgage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MortgageRepository defines methods for mortgage storage operations
type MortgageRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*Mortgage, error)
	FindByFinancialAnalysisID(ctx context.Context, analysisID primitive.ObjectID) (*Mortgage, error)
	Save(ctx context.Context, mortgage *Mortgage) (primitive.ObjectID, error)
	Update(ctx context.Context, mortgage *Mortgage) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Payment schedule methods
	AddMortgageInterval(ctx context.Context, mortgageID primitive.ObjectID, interval *MortgageInterval) error
	UpdateMortgageInterval(ctx context.Context, interval *MortgageInterval) error
	DeleteMortgageInterval(ctx context.Context, intervalID primitive.ObjectID) error
	ClearMortgageIntervals(ctx context.Context, mortgageID primitive.ObjectID) error
}
