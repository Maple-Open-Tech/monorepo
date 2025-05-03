// cloud/backend/internal/ipe/usecase/evaluation/addpropertyphoto.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type AddPropertyPhotoUseCase interface {
	Execute(ctx context.Context, evaluationID primitive.ObjectID, photo *dom_evaluation.PropertyPhoto) error
}

type addPropertyPhotoUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewAddPropertyPhotoUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) AddPropertyPhotoUseCase {
	return &addPropertyPhotoUseCaseImpl{config, logger, repo}
}

func (uc *addPropertyPhotoUseCaseImpl) Execute(ctx context.Context, evaluationID primitive.ObjectID, photo *dom_evaluation.PropertyPhoto) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if evaluationID.IsZero() {
		e["evaluation_id"] = "Evaluation ID is required"
	}
	if photo == nil {
		e["photo"] = "Property photo is required"
	} else {
		if photo.PhotoName == "" {
			e["photo_name"] = "Photo name is required"
		}
		// Additional validations
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating add property photo",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Add to database.
	//
	return uc.repo.AddPropertyPhoto(ctx, evaluationID, photo)
}
