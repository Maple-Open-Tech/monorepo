package evaluation

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EvaluationRepository defines methods for evaluation storage operations
type EvaluationRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*Evaluation, error)
	FindByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*Evaluation, error)
	Save(ctx context.Context, evaluation *Evaluation) (primitive.ObjectID, error)
	Update(ctx context.Context, evaluation *Evaluation) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Building and Legal operations are done through evaluation updates

	// Neighbourhood operations
	FindNeighbourhoodByID(ctx context.Context, id primitive.ObjectID) (*Neighbourhood, error)
	SaveNeighbourhood(ctx context.Context, neighbourhood *Neighbourhood) (primitive.ObjectID, error)
	UpdateNeighbourhood(ctx context.Context, neighbourhood *Neighbourhood) error
	DeleteNeighbourhood(ctx context.Context, id primitive.ObjectID) error

	// Property photo operations
	AddPropertyPhoto(ctx context.Context, evaluationID primitive.ObjectID, photo *PropertyPhoto) error
	UpdatePropertyPhoto(ctx context.Context, photo *PropertyPhoto) error
	DeletePropertyPhoto(ctx context.Context, photoID primitive.ObjectID) error
	FindPhotosByEvaluationID(ctx context.Context, evaluationID primitive.ObjectID) ([]*PropertyPhoto, error)
}
