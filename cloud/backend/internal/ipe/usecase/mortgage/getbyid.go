// cloud/backend/internal/ipe/usecase/mortgage/getbyid.go
package mortgage

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type MortgageGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_mortgage.Mortgage, error)
}

type mortgageGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_mortgage.MortgageRepository
}

func NewMortgageGetByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_mortgage.MortgageRepository,
) MortgageGetByIDUseCase {
	return &mortgageGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *mortgageGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_mortgage.Mortgage, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating mortgage get by ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByID(ctx, id)
}
