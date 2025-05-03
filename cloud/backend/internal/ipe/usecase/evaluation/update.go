// cloud/backend/internal/ipe/usecase/evaluation/update.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type EvaluationUpdateUseCase interface {
	Execute(ctx context.Context, evaluation *dom_evaluation.Evaluation) error
}

type evaluationUpdateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewEvaluationUpdateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) EvaluationUpdateUseCase {
	return &evaluationUpdateUseCaseImpl{config, logger, repo}
}

func (uc *evaluationUpdateUseCaseImpl) Execute(ctx context.Context, evaluation *dom_evaluation.Evaluation) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if evaluation == nil {
		e["evaluation"] = "Evaluation is required"
	} else {
		if evaluation.ID.IsZero() {
			e["id"] = "ID is required"
		}
		// Additional validations
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating evaluation update",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update in database.
	//
	return uc.repo.Update(ctx, evaluation)
}
