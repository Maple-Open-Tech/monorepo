// cloud/backend/internal/ipe/usecase/incomeproperty/delete.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type IncomePropertyDeleteUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) error
}

type incomePropertyDeleteUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyDeleteUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyDeleteUseCase {
	return &incomePropertyDeleteUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyDeleteUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) error {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating property delete",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Delete from database.
	//
	return uc.repo.Delete(ctx, id)
}
