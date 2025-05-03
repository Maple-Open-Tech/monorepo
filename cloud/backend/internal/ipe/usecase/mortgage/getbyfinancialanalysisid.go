// cloud/backend/internal/ipe/usecase/mortgage/getbyfinancialanalysisid.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type MortgageGetByFinancialAnalysisIDUseCase interface {
	Execute(ctx context.Context, analysisID primitive.ObjectID) (*dom_mortgage.Mortgage, error)
}

type mortgageGetByFinancialAnalysisIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_mortgage.MortgageRepository
}

func NewMortgageGetByFinancialAnalysisIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_mortgage.MortgageRepository,
) MortgageGetByFinancialAnalysisIDUseCase {
	return &mortgageGetByFinancialAnalysisIDUseCaseImpl{config, logger, repo}
}

func (uc *mortgageGetByFinancialAnalysisIDUseCaseImpl) Execute(ctx context.Context, analysisID primitive.ObjectID) (*dom_mortgage.Mortgage, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if analysisID.IsZero() {
		e["analysis_id"] = "Financial analysis ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating mortgage get by financial analysis ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByFinancialAnalysisID(ctx, analysisID)
}
