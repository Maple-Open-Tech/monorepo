// cloud/backend/internal/ipe/usecase/person/client/create.go
package client

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type ClientCreateUseCase interface {
	Execute(ctx context.Context, client *dom_person.Client) (primitive.ObjectID, error)
}

type clientCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_person.PersonRepository
}

func NewClientCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_person.PersonRepository,
) ClientCreateUseCase {
	return &clientCreateUseCaseImpl{config, logger, repo}
}

func (uc *clientCreateUseCaseImpl) Execute(ctx context.Context, client *dom_person.Client) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if client == nil {
		e["client"] = "Client is required"
	} else {
		if client.PersonName == "" {
			e["person_name"] = "Person name is required"
		}
		// Additional validations
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating client creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.SaveClient(ctx, client)
}
