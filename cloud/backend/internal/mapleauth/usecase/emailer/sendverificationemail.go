package emailer

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/domain/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/repo/templatedemailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type SendUserVerificationEmailUseCase interface {
	Execute(ctx context.Context, monolithModule int, user *domain.FederatedUser) error
}
type sendUserVerificationEmailUseCaseImpl struct {
	config  *config.Configuration
	logger  *zap.Logger
	emailer templatedemailer.TemplatedEmailer
}

func NewSendUserVerificationEmailUseCase(config *config.Configuration, logger *zap.Logger, emailer templatedemailer.TemplatedEmailer) SendUserVerificationEmailUseCase {
	return &sendUserVerificationEmailUseCaseImpl{config, logger, emailer}
}

func (uc *sendUserVerificationEmailUseCaseImpl) Execute(ctx context.Context, monolithModule int, user *domain.FederatedUser) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if user == nil {
		e["user"] = "FederatedUser is missing value"
	} else {
		if user.FirstName == "" {
			e["first_name"] = "First name is required"
		}
		if user.Email == "" {
			e["email"] = "Email is required"
		}
		if user.EmailVerificationCode == "" {
			e["email_verification_code"] = "Email verification code is required"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			zap.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Send email
	//

	return uc.emailer.SendUserVerificationEmail(ctx, monolithModule, user.Email, user.EmailVerificationCode, user.FirstName)
}
