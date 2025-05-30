package bannedipaddress

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_banip "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateBannedIPAddressUseCase interface {
	Execute(ctx context.Context, bannedIPAddress *dom_banip.BannedIPAddress) error
}

type createBannedIPAddressUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_banip.Repository
}

func NewCreateBannedIPAddressUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_banip.Repository,
) CreateBannedIPAddressUseCase {
	return &createBannedIPAddressUseCaseImpl{config, logger, repo}
}

func (uc *createBannedIPAddressUseCaseImpl) Execute(ctx context.Context, bannedIPAddress *dom_banip.BannedIPAddress) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if bannedIPAddress == nil {
		e["banned_ip_address"] = "Banned IP address is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Create(ctx, bannedIPAddress)
}
