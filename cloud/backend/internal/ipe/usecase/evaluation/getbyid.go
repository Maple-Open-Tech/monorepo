// cloud/backend/internal/ipe/usecase/evaluation/getbyid.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type EvaluationGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_evaluation.Evaluation, error)
}

type evaluationGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewEvaluationGetByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) EvaluationGetByIDUseCase {
	return &evaluationGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *evaluationGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_evaluation.Evaluation, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating evaluation get by ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByID(ctx, id)
}
