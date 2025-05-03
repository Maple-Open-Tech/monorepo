// cloud/backend/internal/ipe/usecase/evaluation/create.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type EvaluationCreateUseCase interface {
	Execute(ctx context.Context, evaluation *dom_evaluation.Evaluation) (primitive.ObjectID, error)
}

type evaluationCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewEvaluationCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) EvaluationCreateUseCase {
	return &evaluationCreateUseCaseImpl{config, logger, repo}
}

func (uc *evaluationCreateUseCaseImpl) Execute(ctx context.Context, evaluation *dom_evaluation.Evaluation) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if evaluation == nil {
		e["evaluation"] = "Evaluation is required"
	} else {
		if evaluation.PropertyID.IsZero() {
			e["property_id"] = "Property ID is required"
		}
		if evaluation.ClientID.IsZero() {
			e["client_id"] = "Client ID is required"
		}
		if evaluation.PresenterID.IsZero() {
			e["presenter_id"] = "Presenter ID is required"
		}
		// Additional validations as needed
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating evaluation creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.Save(ctx, evaluation)
}
