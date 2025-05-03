// cloud/backend/internal/ipe/service/person/client/get.go
package client

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	uc_client "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/client"
)

type ClientResponseDTO struct {
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
	HasLogo         bool               `json:"hasLogo"`
}

type GetClientService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*ClientResponseDTO, error)
}

type getClientServiceImpl struct {
	config               *config.Configuration
	logger               *zap.Logger
	clientGetByIDUseCase uc_client.ClientGetByIDUseCase
}

func NewGetClientService(
	config *config.Configuration,
	logger *zap.Logger,
	clientGetByIDUseCase uc_client.ClientGetByIDUseCase,
) GetClientService {
	return &getClientServiceImpl{
		config:               config,
		logger:               logger,
		clientGetByIDUseCase: clientGetByIDUseCase,
	}
}

func (svc *getClientServiceImpl) Execute(ctx context.Context, id primitive.ObjectID) (*ClientResponseDTO, error) {
	// Get client
	client, err := svc.clientGetByIDUseCase.Execute(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Client not found", zap.String("id", id.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get client", zap.Error(err))
		return nil, err
	}
	if client == nil {
		svc.logger.Warn("Client not found", zap.String("id", id.Hex()))
		return nil, errors.New("client not found")
	}

	// Create response
	response := &ClientResponseDTO{
		ID:              client.ID,
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
		HasLogo:         len(client.LogoPhotoData) > 0,
	}

	return response, nil
}
