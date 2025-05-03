// cloud/backend/internal/ipe/service/incomeproperty/create.go
package incomeproperty

import (
	"context"
	"time"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
	uc_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateIncomePropertyRequestDTO struct {
	Address      string `json:"address"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
	PropertyCode string `json:"propertyCode"`
	RecordName   string `json:"recordName"`
}

type CreateIncomePropertyResponseDTO struct {
	ID           primitive.ObjectID `json:"id"`
	Address      string             `json:"address"`
	City         string             `json:"city"`
	Province     string             `json:"province"`
	Country      string             `json:"country"`
	PropertyCode string             `json:"propertyCode"`
	RecordName   string             `json:"recordName"`
}

type CreateIncomePropertyService interface {
	Execute(ctx context.Context, request *CreateIncomePropertyRequestDTO) (*CreateIncomePropertyResponseDTO, error)
}

type createIncomePropertyServiceImpl struct {
	config                      *config.Configuration
	logger                      *zap.Logger
	incomePropertyCreateUseCase uc_property.IncomePropertyCreateUseCase
}

func NewCreateIncomePropertyService(
	config *config.Configuration,
	logger *zap.Logger,
	incomePropertyCreateUseCase uc_property.IncomePropertyCreateUseCase,
) CreateIncomePropertyService {
	return &createIncomePropertyServiceImpl{
		config:                      config,
		logger:                      logger,
		incomePropertyCreateUseCase: incomePropertyCreateUseCase,
	}
}

func (svc *createIncomePropertyServiceImpl) Execute(ctx context.Context, req *CreateIncomePropertyRequestDTO) (*CreateIncomePropertyResponseDTO, error) {
	// Validate request
	if req == nil {
		return nil, httperror.NewForBadRequestWithSingleField("request", "Request is required")
	}

	errors := make(map[string]string)
	if req.Address == "" {
		errors["address"] = "Address is required"
	}
	if req.City == "" {
		errors["city"] = "City is required"
	}
	if req.Province == "" {
		errors["province"] = "Province is required"
	}
	if req.Country == "" {
		errors["country"] = "Country is required"
	}
	if req.PropertyCode == "" {
		errors["propertyCode"] = "Property code is required"
	}

	if len(errors) > 0 {
		return nil, httperror.NewForBadRequest(&errors)
	}

	// Create property domain object
	property := &dom_property.IncomeProperty{
		ID:                 primitive.NewObjectID(),
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		Country:            req.Country,
		PropertyCode:       req.PropertyCode,
		RecordName:         req.RecordName,
		RecordCreationDate: time.Now(),
	}

	// If record name is empty, use address
	if property.RecordName == "" {
		property.RecordName = property.Address
	}

	// Create property
	id, err := svc.incomePropertyCreateUseCase.Execute(ctx, property)
	if err != nil {
		svc.logger.Error("Failed to create property", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateIncomePropertyResponseDTO{
		ID:           id,
		Address:      property.Address,
		City:         property.City,
		Province:     property.Province,
		Country:      property.Country,
		PropertyCode: property.PropertyCode,
		RecordName:   property.RecordName,
	}

	return response, nil
}
