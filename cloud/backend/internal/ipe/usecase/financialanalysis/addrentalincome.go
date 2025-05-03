// cloud/backend/internal/ipe/usecase/financialanalysis/addrentalincome.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type AddRentalIncomeUseCase interface {
	Execute(ctx context.Context, analysisID primitive.ObjectID, income *dom_financial.RentalIncome) error
}

type addRentalIncomeUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_financial.FinancialRepository
}

func NewAddRentalIncomeUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_financial.FinancialRepository,
) AddRentalIncomeUseCase {
	return &addRentalIncomeUseCaseImpl{config, logger, repo}
}

func (uc *addRentalIncomeUseCaseImpl) Execute(ctx context.Context, analysisID primitive.ObjectID, income *dom_financial.RentalIncome) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if analysisID.IsZero() {
		e["analysis_id"] = "Analysis ID is required"
	}
	if income == nil {
		e["income"] = "Rental income is required"
	} else {
		if income.NameText == "" {
			e["name_text"] = "Name is required"
		}
		// Additional validations as needed
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating add rental income",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Add to database.
	//
	return uc.repo.AddRentalIncome(ctx, analysisID, income)
}
