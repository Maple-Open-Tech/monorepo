// cloud/backend/internal/ipe/usecase/incomeproperty/update.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type IncomePropertyUpdateUseCase interface {
	Execute(ctx context.Context, property *dom_property.IncomeProperty) error
}

type incomePropertyUpdateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyUpdateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyUpdateUseCase {
	return &incomePropertyUpdateUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyUpdateUseCaseImpl) Execute(ctx context.Context, property *dom_property.IncomeProperty) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if property == nil {
		e["property"] = "Property is required"
	} else {
		if property.ID.IsZero() {
			e["id"] = "ID is required"
		}
		if property.Address == "" {
			e["address"] = "Address is required"
		}
		if property.City == "" {
			e["city"] = "City is required"
		}
		if property.Province == "" {
			e["province"] = "Province is required"
		}
		if property.Country == "" {
			e["country"] = "Country is required"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating property update",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update in database.
	//
	return uc.repo.Update(ctx, property)
}
