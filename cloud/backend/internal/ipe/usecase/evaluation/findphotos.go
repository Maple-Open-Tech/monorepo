// cloud/backend/internal/ipe/usecase/evaluation/findphotos.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FindPhotosByEvaluationIDUseCase interface {
	Execute(ctx context.Context, evaluationID primitive.ObjectID) ([]*dom_evaluation.PropertyPhoto, error)
}

type findPhotosByEvaluationIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewFindPhotosByEvaluationIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) FindPhotosByEvaluationIDUseCase {
	return &findPhotosByEvaluationIDUseCaseImpl{config, logger, repo}
}

func (uc *findPhotosByEvaluationIDUseCaseImpl) Execute(ctx context.Context, evaluationID primitive.ObjectID) ([]*dom_evaluation.PropertyPhoto, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if evaluationID.IsZero() {
		e["evaluation_id"] = "Evaluation ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating find photos by evaluation ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindPhotosByEvaluationID(ctx, evaluationID)
}
