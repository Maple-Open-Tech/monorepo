// cloud/backend/internal/ipe/usecase/mortgage/addmortgageinterval.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type AddMortgageIntervalUseCase interface {
	Execute(ctx context.Context, mortgageID primitive.ObjectID, interval *dom_mortgage.MortgageInterval) error
}

type addMortgageIntervalUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_mortgage.MortgageRepository
}

func NewAddMortgageIntervalUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_mortgage.MortgageRepository,
) AddMortgageIntervalUseCase {
	return &addMortgageIntervalUseCaseImpl{config, logger, repo}
}

func (uc *addMortgageIntervalUseCaseImpl) Execute(ctx context.Context, mortgageID primitive.ObjectID, interval *dom_mortgage.MortgageInterval) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if mortgageID.IsZero() {
		e["mortgage_id"] = "Mortgage ID is required"
	}
	if interval == nil {
		e["interval"] = "Mortgage interval is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating add mortgage interval",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Add to database.
	//
	return uc.repo.AddMortgageInterval(ctx, mortgageID, interval)
}
