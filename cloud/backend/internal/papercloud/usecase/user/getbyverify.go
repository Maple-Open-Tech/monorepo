package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type UserGetByVerificationCodeUseCase interface {
	Execute(ctx context.Context, verificationCode string) (*dom_user.User, error)
}

type userGetByVerificationCodeUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_user.Repository
}

func NewUserGetByVerificationCodeUseCase(config *config.Configuration, logger *zap.Logger, repo dom_user.Repository) UserGetByVerificationCodeUseCase {
	return &userGetByVerificationCodeUseCaseImpl{config, logger, repo}
}

func (uc *userGetByVerificationCodeUseCaseImpl) Execute(ctx context.Context, verificationCode string) (*dom_user.User, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if verificationCode == "" {
		e["verification_code"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for get by verification",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Get from database.
	//

	return uc.repo.GetByVerificationCode(ctx, verificationCode)
}
