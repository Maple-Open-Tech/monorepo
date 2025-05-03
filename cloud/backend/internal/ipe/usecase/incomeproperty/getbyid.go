// cloud/backend/internal/ipe/usecase/incomeproperty/getbyid.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type IncomePropertyGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_property.IncomeProperty, error)
}

type incomePropertyGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyGetByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyGetByIDUseCase {
	return &incomePropertyGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_property.IncomeProperty, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating property get by ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindByID(ctx, id)
}
