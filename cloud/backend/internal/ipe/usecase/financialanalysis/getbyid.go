// cloud/backend/internal/ipe/usecase/financialanalysis/getbyid.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FinancialAnalysisGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_financial.FinancialAnalysis, error)
}

type financialAnalysisGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_financial.FinancialRepository
}

func NewFinancialAnalysisGetByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_financial.FinancialRepository,
) FinancialAnalysisGetByIDUseCase {
	return &financialAnalysisGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *financialAnalysisGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_financial.FinancialAnalysis, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating financial analysis get by ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByID(ctx, id)
}
