// cloud/backend/internal/ipe/usecase/incomeproperty/findbycity.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type IncomePropertyFindByCityUseCase interface {
	Execute(ctx context.Context, city string) ([]dom_property.IncomeProperty, error)
}

type incomePropertyFindByCityUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyFindByCityUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyFindByCityUseCase {
	return &incomePropertyFindByCityUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyFindByCityUseCaseImpl) Execute(ctx context.Context, city string) ([]dom_property.IncomeProperty, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if city == "" {
		e["city"] = "City is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating property find by city",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Find in database.
	//
	return uc.repo.FindByCity(ctx, city)
}
