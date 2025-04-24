package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/emailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewCreateBannedIPAddressUseCase,
			bannedipaddress.NewBannedIPAddressListAllValuesUseCase,
			emailer.NewSendUserPasswordResetEmailUseCase,
			emailer.NewSendUserVerificationEmailUseCase,
			user.NewUserGetBySessionIDUseCase,
			user.NewUserCountByFilterUseCase,
			user.NewUserCreateUseCase,
			user.NewUserDeleteUserByEmailUseCase,
			user.NewUserDeleteByIDUseCase,
			user.NewUserGetByEmailUseCase,
			user.NewUserGetByIDUseCase,
			user.NewUserGetByVerificationCodeUseCase,
			user.NewUserListAllUseCase,
			user.NewUserListByFilterUseCase,
			user.NewUserUpdateUseCase,
		),
	)
}
