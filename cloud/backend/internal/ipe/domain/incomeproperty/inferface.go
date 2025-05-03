package incomeproperty

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PropertyRepository defines methods for property storage operations
type PropertyRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*IncomeProperty, error)
	FindAll(ctx context.Context) ([]*IncomeProperty, error)
	Save(ctx context.Context, property *IncomeProperty) (primitive.ObjectID, error)
	Update(ctx context.Context, property *IncomeProperty) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByAddress(ctx context.Context, address string) ([]IncomeProperty, error)
	FindByCity(ctx context.Context, city string) ([]IncomeProperty, error)
}
