// cloud/backend/internal/ipe/usecase/financialanalysis/create.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FinancialAnalysisCreateUseCase interface {
	Execute(ctx context.Context, analysis *dom_financial.FinancialAnalysis) (primitive.ObjectID, error)
}

type financialAnalysisCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_financial.FinancialRepository
}

func NewFinancialAnalysisCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_financial.FinancialRepository,
) FinancialAnalysisCreateUseCase {
	return &financialAnalysisCreateUseCaseImpl{config, logger, repo}
}

func (uc *financialAnalysisCreateUseCaseImpl) Execute(ctx context.Context, analysis *dom_financial.FinancialAnalysis) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if analysis == nil {
		e["analysis"] = "Financial analysis is required"
	} else {
		if analysis.PropertyID.IsZero() {
			e["property_id"] = "Property ID is required"
		}
		// Additional validations can be added as needed
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating financial analysis creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.Save(ctx, analysis)
}
