// cloud/backend/internal/ipe/usecase/financialanalysis/getbypropertyid.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FinancialAnalysisGetByPropertyIDUseCase interface {
	Execute(ctx context.Context, propertyID primitive.ObjectID) (*dom_financial.FinancialAnalysis, error)
}

type financialAnalysisGetByPropertyIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_financial.FinancialRepository
}

func NewFinancialAnalysisGetByPropertyIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_financial.FinancialRepository,
) FinancialAnalysisGetByPropertyIDUseCase {
	return &financialAnalysisGetByPropertyIDUseCaseImpl{config, logger, repo}
}

func (uc *financialAnalysisGetByPropertyIDUseCaseImpl) Execute(ctx context.Context, propertyID primitive.ObjectID) (*dom_financial.FinancialAnalysis, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if propertyID.IsZero() {
		e["property_id"] = "Property ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating financial analysis get by property ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByPropertyID(ctx, propertyID)
}
