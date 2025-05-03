// cloud/backend/internal/ipe/service/incomeproperty/list.go
package incomeproperty

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	uc_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/incomeproperty"
)

type ListIncomePropertiesResponseDTO struct {
	Properties []IncomePropertyResponseDTO `json:"properties"`
}

type ListIncomePropertiesService interface {
	Execute(ctx context.Context) (*ListIncomePropertiesResponseDTO, error)
}

type listIncomePropertiesServiceImpl struct {
	config                       *config.Configuration
	logger                       *zap.Logger
	incomePropertyListAllUseCase uc_property.IncomePropertyListAllUseCase
}

func NewListIncomePropertiesService(
	config *config.Configuration,
	logger *zap.Logger,
	incomePropertyListAllUseCase uc_property.IncomePropertyListAllUseCase,
) ListIncomePropertiesService {
	return &listIncomePropertiesServiceImpl{
		config:                       config,
		logger:                       logger,
		incomePropertyListAllUseCase: incomePropertyListAllUseCase,
	}
}

func (svc *listIncomePropertiesServiceImpl) Execute(ctx context.Context) (*ListIncomePropertiesResponseDTO, error) {
	// Get all properties
	properties, err := svc.incomePropertyListAllUseCase.Execute(ctx)
	if err != nil {
		svc.logger.Error("Failed to list properties", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &ListIncomePropertiesResponseDTO{
		Properties: make([]IncomePropertyResponseDTO, len(properties)),
	}

	// Map domain objects to DTOs
	for i, property := range properties {
		response.Properties[i] = mapPropertyToDTO(property)
	}

	return response, nil
}

func mapPropertyToDTO(property *dom_property.IncomeProperty) IncomePropertyResponseDTO {
	return IncomePropertyResponseDTO{
		ID:                 property.ID,
		Address:            property.Address,
		City:               property.City,
		Province:           property.Province,
		Country:            property.Country,
		PropertyCode:       property.PropertyCode,
		RecordName:         property.RecordName,
		RecordCreationDate: property.RecordCreationDate.Format("2006-01-02"),
		HasMainPhoto:       len(property.MainPhotoThumbnail) > 0,
	}
}
