// cloud/backend/internal/ipe/usecase/evaluation/getbypropertyid.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type EvaluationGetByPropertyIDUseCase interface {
	Execute(ctx context.Context, propertyID primitive.ObjectID) (*dom_evaluation.Evaluation, error)
}

type evaluationGetByPropertyIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_evaluation.EvaluationRepository
}

func NewEvaluationGetByPropertyIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_evaluation.EvaluationRepository,
) EvaluationGetByPropertyIDUseCase {
	return &evaluationGetByPropertyIDUseCaseImpl{config, logger, repo}
}

func (uc *evaluationGetByPropertyIDUseCaseImpl) Execute(ctx context.Context, propertyID primitive.ObjectID) (*dom_evaluation.Evaluation, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if propertyID.IsZero() {
		e["property_id"] = "Property ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating evaluation get by property ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByPropertyID(ctx, propertyID)
}
