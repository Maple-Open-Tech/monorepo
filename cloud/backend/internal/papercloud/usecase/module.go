package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/usecase/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/usecase/emailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/usecase/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewCreateBannedIPAddressUseCase,
			bannedipaddress.NewBannedIPAddressListAllValuesUseCase,
			emailer.NewSendUserPasswordResetEmailUseCase,
			emailer.NewSendUserVerificationEmailUseCase,
			user.NewUserCountByFilterUseCase,
			user.NewUserCreateUseCase,
			user.NewUserUpdateUseCase,
			user.NewUserListByFilterUseCase,
			user.NewUserListAllUseCase,
			user.NewUserGetByVerificationCodeUseCase,
			user.NewUserGetBySessionIDUseCase,
			user.NewUserGetByIDUseCase,
			user.NewUserGetByEmailUseCase,
			user.NewUserDeleteByIDUseCase,
			user.NewUserDeleteUserByEmailUseCase,
		),
	)
}
