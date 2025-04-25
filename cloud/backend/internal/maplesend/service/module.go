package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service/me"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			gateway.NewGatewayUserRegisterService,
			gateway.NewGatewayLoginService,
			gateway.NewGatewayLogoutService,
			gateway.NewGatewayResetPasswordService,
			gateway.NewGatewaySendVerifyEmailService,
			gateway.NewGatewayVerifyEmailService,
			gateway.NewGatewayRefreshTokenService,
			gateway.NewGatewayForgotPasswordService,
			me.NewGetMeService,
			me.NewUpdateMeService,
			me.NewVerifyProfileService,
			me.NewDeleteMeService,
		),
	)
}
