// cloud/backend/internal/ipe/service/incomeproperty/get.go
package incomeproperty

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	uc_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/incomeproperty"
)

type IncomePropertyResponseDTO struct {
	ID                 primitive.ObjectID `json:"id"`
	Address            string             `json:"address"`
	City               string             `json:"city"`
	Province           string             `json:"province"`
	Country            string             `json:"country"`
	PropertyCode       string             `json:"propertyCode"`
	RecordName         string             `json:"recordName"`
	RecordCreationDate string             `json:"recordCreationDate"`
	HasMainPhoto       bool               `json:"hasMainPhoto"`
}

type GetIncomePropertyService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*IncomePropertyResponseDTO, error)
}

type getIncomePropertyServiceImpl struct {
	config                 *config.Configuration
	logger                 *zap.Logger
	propertyGetByIDUseCase uc_property.IncomePropertyGetByIDUseCase
}

func NewGetIncomePropertyService(
	config *config.Configuration,
	logger *zap.Logger,
	propertyGetByIDUseCase uc_property.IncomePropertyGetByIDUseCase,
) GetIncomePropertyService {
	return &getIncomePropertyServiceImpl{
		config:                 config,
		logger:                 logger,
		propertyGetByIDUseCase: propertyGetByIDUseCase,
	}
}

func (svc *getIncomePropertyServiceImpl) Execute(ctx context.Context, id primitive.ObjectID) (*IncomePropertyResponseDTO, error) {
	// Get property
	property, err := svc.propertyGetByIDUseCase.Execute(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Property not found", zap.String("id", id.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get property", zap.Error(err))
		return nil, err
	}
	if property == nil {
		svc.logger.Warn("Property not found", zap.String("id", id.Hex()))
		return nil, errors.New("property not found")
	}

	// Create response
	response := &IncomePropertyResponseDTO{
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

	return response, nil
}
