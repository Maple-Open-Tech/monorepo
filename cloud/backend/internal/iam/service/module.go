package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service/token"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			token.NewTokenVerifyService,
			token.NewTokenGetSessionService,
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
