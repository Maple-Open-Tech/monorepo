package bannedipaddress

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_banip "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/bannedipaddress"
)

type BannedIPAddressListAllValuesUseCase interface {
	Execute(ctx context.Context) ([]string, error)
}

type bannedIPAddressListAllValuesUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_banip.Repository
}

func NewBannedIPAddressListAllValuesUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_banip.Repository,
) BannedIPAddressListAllValuesUseCase {
	return &bannedIPAddressListAllValuesUseCaseImpl{config, logger, repo}
}

func (uc *bannedIPAddressListAllValuesUseCaseImpl) Execute(ctx context.Context) ([]string, error) {
	return uc.repo.ListAllValues(ctx)
}
