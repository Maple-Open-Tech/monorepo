// cloud/backend/internal/ipe/usecase/incomeproperty/create.go
package incomeproperty

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type IncomePropertyCreateUseCase interface {
	Execute(ctx context.Context, property *dom_property.IncomeProperty) (primitive.ObjectID, error)
}

type incomePropertyCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_property.PropertyRepository
}

func NewIncomePropertyCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_property.PropertyRepository,
) IncomePropertyCreateUseCase {
	return &incomePropertyCreateUseCaseImpl{config, logger, repo}
}

func (uc *incomePropertyCreateUseCaseImpl) Execute(ctx context.Context, property *dom_property.IncomeProperty) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if property == nil {
		e["property"] = "Property is required"
	} else {
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
		if property.PropertyCode == "" {
			e["property_code"] = "Property code is required"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating property creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	if property.RecordName == "" {
		property.RecordName = property.Address
	}
	if property.RecordCreationDate.IsZero() {
		property.RecordCreationDate = time.Now()
	}

	return uc.repo.Save(ctx, property)
}
