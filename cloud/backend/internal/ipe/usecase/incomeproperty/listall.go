// cloud/backend/internal/ipe/usecase/incomeproperty/listall.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
)

type IncomePropertyListAllUseCase interface {
	Execute(ctx context.Context) ([]*dom_property.IncomeProperty, error)
}

type incomePropertyListAllUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyListAllUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyListAllUseCase {
	return &incomePropertyListAllUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyListAllUseCaseImpl) Execute(ctx context.Context) ([]*dom_property.IncomeProperty, error) {
	uc.logger.Debug("executing list all properties use case")

	properties, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("failed to list all properties", zap.Any("error", err))
		return nil, err
	}

	return properties, nil
}
