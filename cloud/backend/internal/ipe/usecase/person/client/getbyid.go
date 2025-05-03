// cloud/backend/internal/ipe/usecase/person/client/getbyid.go
package client

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type ClientGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_person.Client, error)
}

type clientGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_person.PersonRepository
}

func NewClientGetByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_person.PersonRepository,
) ClientGetByIDUseCase {
	return &clientGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *clientGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_person.Client, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating client get by ID",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//
	return uc.repo.FindClientByID(ctx, id)
}
