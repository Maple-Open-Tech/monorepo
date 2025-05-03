// cloud/backend/internal/ipe/service/person/owner/create.go
package owner

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	uc_owner "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/owner"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateOwnerRequestDTO struct {
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

type CreateOwnerResponseDTO struct {
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

type CreateOwnerService interface {
	Execute(ctx context.Context, request *CreateOwnerRequestDTO) (*CreateOwnerResponseDTO, error)
}

type createOwnerServiceImpl struct {
	config             *config.Configuration
	logger             *zap.Logger
	ownerCreateUseCase uc_owner.OwnerCreateUseCase
}

func NewCreateOwnerService(
	config *config.Configuration,
	logger *zap.Logger,
	ownerCreateUseCase uc_owner.OwnerCreateUseCase,
) CreateOwnerService {
	return &createOwnerServiceImpl{
		config:             config,
		logger:             logger,
		ownerCreateUseCase: ownerCreateUseCase,
	}
}

func (svc *createOwnerServiceImpl) Execute(ctx context.Context, req *CreateOwnerRequestDTO) (*CreateOwnerResponseDTO, error) {
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

	// Create owner domain object
	owner := &dom_person.Owner{
		ID: primitive.NewObjectID(),
	}

	// Set fields through assignment
	owner.PersonName = req.PersonName
	owner.Address = req.Address
	owner.City = req.City
	owner.Province = req.Province
	owner.Country = req.Country
	owner.PostalCode = req.PostalCode
	owner.Email = req.Email
	owner.OfficeTelNumber = req.OfficeTelNumber
	owner.MobileTelNumber = req.MobileTelNumber
	owner.FaxTelNumber = req.FaxTelNumber
	owner.Website = req.Website
	owner.RecordUniqueID = req.RecordUniqueID

	// Create record unique ID if not provided
	if owner.RecordUniqueID == "" {
		owner.RecordUniqueID = primitive.NewObjectID().Hex()
	}

	// Create owner
	id, err := svc.ownerCreateUseCase.Execute(ctx, owner)
	if err != nil {
		svc.logger.Error("Failed to create owner", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateOwnerResponseDTO{
		ID:              id,
		PersonName:      owner.PersonName,
		Address:         owner.Address,
		City:            owner.City,
		Province:        owner.Province,
		Country:         owner.Country,
		PostalCode:      owner.PostalCode,
		Email:           owner.Email,
		OfficeTelNumber: owner.OfficeTelNumber,
		MobileTelNumber: owner.MobileTelNumber,
		FaxTelNumber:    owner.FaxTelNumber,
		Website:         owner.Website,
		RecordUniqueID:  owner.RecordUniqueID,
	}

	return response, nil
}
