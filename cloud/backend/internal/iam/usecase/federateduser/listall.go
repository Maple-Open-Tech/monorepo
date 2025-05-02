package federateduser

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
)

type UserListAllUseCase interface {
	Execute(ctx context.Context) ([]*dom_user.FederatedUser, error)
}

type userListAllUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_user.Repository
}

func NewUserListAllUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_user.Repository,
) UserListAllUseCase {
	return &userListAllUseCaseImpl{
		config: config,
		logger: logger,
		repo:   repo,
	}
}

func (uc *userListAllUseCaseImpl) Execute(ctx context.Context) ([]*dom_user.FederatedUser, error) {
	uc.logger.Debug("executing list all users use case")

	users, err := uc.repo.ListAll(ctx)
	if err != nil {
		uc.logger.Error("failed to list all users", zap.Any("error", err))
		return nil, err
	}

	return users, nil
}
