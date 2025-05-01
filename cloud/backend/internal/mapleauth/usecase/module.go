package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/baseuser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/emailer"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewCreateBannedIPAddressUseCase,
			bannedipaddress.NewBannedIPAddressListAllValuesUseCase,
			emailer.NewSendUserPasswordResetEmailUseCase,
			emailer.NewSendUserVerificationEmailUseCase,
			baseuser.NewUserGetBySessionIDUseCase,
			baseuser.NewUserCountByFilterUseCase,
			baseuser.NewUserCreateUseCase,
			baseuser.NewUserDeleteUserByEmailUseCase,
			baseuser.NewUserDeleteByIDUseCase,
			baseuser.NewUserGetByEmailUseCase,
			baseuser.NewUserGetByIDUseCase,
			baseuser.NewUserGetByVerificationCodeUseCase,
			baseuser.NewUserListAllUseCase,
			baseuser.NewUserListByFilterUseCase,
			baseuser.NewUserUpdateUseCase,
		),
	)
}
