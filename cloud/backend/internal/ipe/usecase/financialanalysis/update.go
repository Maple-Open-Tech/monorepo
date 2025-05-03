// cloud/backend/internal/ipe/usecase/financialanalysis/update.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FinancialAnalysisUpdateUseCase interface {
	Execute(ctx context.Context, analysis *dom_financial.FinancialAnalysis) error
}

type financialAnalysisUpdateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_financial.FinancialRepository
}

func NewFinancialAnalysisUpdateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_financial.FinancialRepository,
) FinancialAnalysisUpdateUseCase {
	return &financialAnalysisUpdateUseCaseImpl{config, logger, repo}
}

func (uc *financialAnalysisUpdateUseCaseImpl) Execute(ctx context.Context, analysis *dom_financial.FinancialAnalysis) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if analysis == nil {
		e["analysis"] = "Financial analysis is required"
	} else {
		if analysis.ID.IsZero() {
			e["id"] = "ID is required"
		}
		if analysis.PropertyID.IsZero() {
			e["property_id"] = "Property ID is required"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating financial analysis update",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update in database.
	//
	return uc.repo.Update(ctx, analysis)
}
