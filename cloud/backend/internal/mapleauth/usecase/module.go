package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/emailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/federateduser"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewCreateBannedIPAddressUseCase,
			bannedipaddress.NewBannedIPAddressListAllValuesUseCase,
			emailer.NewSendUserPasswordResetEmailUseCase,
			emailer.NewSendUserVerificationEmailUseCase,
			federateduser.NewUserGetBySessionIDUseCase,
			federateduser.NewUserCountByFilterUseCase,
			federateduser.NewUserCreateUseCase,
			federateduser.NewUserDeleteUserByEmailUseCase,
			federateduser.NewUserDeleteByIDUseCase,
			federateduser.NewUserGetByEmailUseCase,
			federateduser.NewUserGetByIDUseCase,
			federateduser.NewUserGetByVerificationCodeUseCase,
			federateduser.NewUserListAllUseCase,
			federateduser.NewUserListByFilterUseCase,
			federateduser.NewUserUpdateUseCase,
		),
	)
}
