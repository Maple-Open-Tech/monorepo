package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service/token"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			token.NewTokenVerifyService,
			token.NewTokenGetSessionService,
			gateway.NewGatewayFederatedUserRegisterService,
			gateway.NewGatewayVerifyEmailService,
			// gateway.NewGatewayLoginService,
			// gateway.NewGatewayLogoutService,
			// gateway.NewGatewaySendVerifyEmailService,
			// gateway.NewGatewayRefreshTokenService,
			// gateway.NewGatewayResetPasswordService,
			// gateway.NewGatewayForgotPasswordService,
			// me.NewGetMeService,
			// me.NewUpdateMeService,
			// me.NewVerifyProfileService,
			// me.NewDeleteMeService,
		),
	)
}
