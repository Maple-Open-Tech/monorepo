// cloud/backend/internal/ipe/service/person/client/create.go
package client

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	uc_client "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/client"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateClientRequestDTO struct {
	PersonName      string `json:"personName"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	PostalCode      string `json:"postalCode"`
	Email           string `json:"email"`
	OfficeTelNumber string `json:"officeTelNumber"`
	MobileTelNumber string `json:"mobileTelNumber"`
	FaxTelNumber    string `json:"faxTelNumber"`
	Website         string `json:"website"`
	RecordUniqueID  string `json:"recordUniqueId"`
}

type CreateClientResponseDTO struct {
	ID              primitive.ObjectID `json:"id"`
	PersonName      string             `json:"personName"`
	Address         string             `json:"address"`
	City            string             `json:"city"`
	Province        string             `json:"province"`
	Country         string             `json:"country"`
	PostalCode      string             `json:"postalCode"`
	Email           string             `json:"email"`
	OfficeTelNumber string             `json:"officeTelNumber"`
	MobileTelNumber string             `json:"mobileTelNumber"`
	FaxTelNumber    string             `json:"faxTelNumber"`
	Website         string             `json:"website"`
	RecordUniqueID  string             `json:"recordUniqueId"`
}

type CreateClientService interface {
	Execute(ctx context.Context, request *CreateClientRequestDTO) (*CreateClientResponseDTO, error)
}

type createClientServiceImpl struct {
	config              *config.Configuration
	logger              *zap.Logger
	clientCreateUseCase uc_client.ClientCreateUseCase
}

func NewCreateClientService(
	config *config.Configuration,
	logger *zap.Logger,
	clientCreateUseCase uc_client.ClientCreateUseCase,
) CreateClientService {
	return &createClientServiceImpl{
		config:              config,
		logger:              logger,
		clientCreateUseCase: clientCreateUseCase,
	}
}

func (svc *createClientServiceImpl) Execute(ctx context.Context, req *CreateClientRequestDTO) (*CreateClientResponseDTO, error) {
	// Validate request
	if req == nil {
		return nil, httperror.NewForBadRequestWithSingleField("request", "Request is required")
	}

	errors := make(map[string]string)
	if req.PersonName == "" {
		errors["personName"] = "Person name is required"
	}
	if req.Email == "" {
		errors["email"] = "Email is required"
	}
	// Additional validations as needed

	if len(errors) > 0 {
		return nil, httperror.NewForBadRequest(&errors)
	}

	// Create client domain object
	client := &dom_person.Client{
		ID: primitive.NewObjectID(),
	}
	// Set fields through assignment
	client.PersonName = req.PersonName
	client.Address = req.Address
	client.City = req.City
	client.Province = req.Province
	client.Country = req.Country
	client.PostalCode = req.PostalCode
	client.Email = req.Email
	client.OfficeTelNumber = req.OfficeTelNumber
	client.MobileTelNumber = req.MobileTelNumber
	client.FaxTelNumber = req.FaxTelNumber
	client.Website = req.Website
	client.RecordUniqueID = req.RecordUniqueID

	// Create record unique ID if not provided
	if client.RecordUniqueID == "" {
		client.RecordUniqueID = primitive.NewObjectID().Hex()
	}

	// Create client
	id, err := svc.clientCreateUseCase.Execute(ctx, client)
	if err != nil {
		svc.logger.Error("Failed to create client", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateClientResponseDTO{
		ID:              id,
		PersonName:      client.PersonName,
		Address:         client.Address,
		City:            client.City,
		Province:        client.Province,
		Country:         client.Country,
		PostalCode:      client.PostalCode,
		Email:           client.Email,
		OfficeTelNumber: client.OfficeTelNumber,
		MobileTelNumber: client.MobileTelNumber,
		FaxTelNumber:    client.FaxTelNumber,
		Website:         client.Website,
		RecordUniqueID:  client.RecordUniqueID,
	}

	return response, nil
}
