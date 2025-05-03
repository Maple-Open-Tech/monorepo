// cloud/backend/internal/ipe/usecase/mortgage/create.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type MortgageCreateUseCase interface {
	Execute(ctx context.Context, mortgage *dom_mortgage.Mortgage) (primitive.ObjectID, error)
}

type mortgageCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_mortgage.MortgageRepository
}

func NewMortgageCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_mortgage.MortgageRepository,
) MortgageCreateUseCase {
	return &mortgageCreateUseCaseImpl{config, logger, repo}
}

func (uc *mortgageCreateUseCaseImpl) Execute(ctx context.Context, mortgage *dom_mortgage.Mortgage) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if mortgage == nil {
		e["mortgage"] = "Mortgage is required"
	} else {
		if mortgage.FinancialAnalysisID.IsZero() {
			e["financial_analysis_id"] = "Financial analysis ID is required"
		}
		// Additional validations can be added as needed
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating mortgage creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.Save(ctx, mortgage)
}
