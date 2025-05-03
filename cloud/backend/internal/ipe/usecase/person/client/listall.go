// cloud/backend/internal/ipe/usecase/person/client/listall.go
package client

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

type ClientListAllUseCase interface {
	Execute(ctx context.Context) ([]*dom_person.Client, error)
}

type clientListAllUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_person.PersonRepository
}

func NewClientListAllUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_person.PersonRepository,
) ClientListAllUseCase {
	return &clientListAllUseCaseImpl{config, logger, repo}
}

func (uc *clientListAllUseCaseImpl) Execute(ctx context.Context) ([]*dom_person.Client, error) {
	uc.logger.Debug("executing list all clients use case")

	clients, err := uc.repo.FindAllClients(ctx)
	if err != nil {
		uc.logger.Error("failed to list all clients", zap.Any("error", err))
		return nil, err
	}

	return clients, nil
}
