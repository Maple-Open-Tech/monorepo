// cloud/backend/internal/ipe/service/person/presenter/create.go
package presenter

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	uc_presenter "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/presenter"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreatePresenterRequestDTO struct {
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

type CreatePresenterResponseDTO struct {
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

type CreatePresenterService interface {
	Execute(ctx context.Context, request *CreatePresenterRequestDTO) (*CreatePresenterResponseDTO, error)
}

type createPresenterServiceImpl struct {
	config                 *config.Configuration
	logger                 *zap.Logger
	presenterCreateUseCase uc_presenter.PresenterCreateUseCase
}

func NewCreatePresenterService(
	config *config.Configuration,
	logger *zap.Logger,
	presenterCreateUseCase uc_presenter.PresenterCreateUseCase,
) CreatePresenterService {
	return &createPresenterServiceImpl{
		config:                 config,
		logger:                 logger,
		presenterCreateUseCase: presenterCreateUseCase,
	}
}

func (svc *createPresenterServiceImpl) Execute(ctx context.Context, req *CreatePresenterRequestDTO) (*CreatePresenterResponseDTO, error) {
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

	// Create presenter domain object
	presenter := &dom_person.Presenter{
		ID: primitive.NewObjectID(),
		// Cannot access basePerson directly
	}
	// We need to create a new presenter and then set the fields on it
	if req.RecordUniqueID == "" {
		req.RecordUniqueID = primitive.NewObjectID().Hex()
	}

	// Create presenter
	id, err := svc.presenterCreateUseCase.Execute(ctx, presenter)
	if err != nil {
		svc.logger.Error("Failed to create presenter", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreatePresenterResponseDTO{
		ID:              id,
		PersonName:      presenter.PersonName,
		Address:         presenter.Address,
		City:            presenter.City,
		Province:        presenter.Province,
		Country:         presenter.Country,
		PostalCode:      presenter.PostalCode,
		Email:           presenter.Email,
		OfficeTelNumber: presenter.OfficeTelNumber,
		MobileTelNumber: presenter.MobileTelNumber,
		FaxTelNumber:    presenter.FaxTelNumber,
		Website:         presenter.Website,
		RecordUniqueID:  presenter.RecordUniqueID,
	}

	return response, nil
}
