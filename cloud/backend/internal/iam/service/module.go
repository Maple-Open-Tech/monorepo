// cloud/backend/internal/iam/service/module.go
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
			// Add the new E2EE login services
			gateway.NewGatewayRequestLoginOTTService,
			gateway.NewGatewayVerifyLoginOTTService,
			gateway.NewGatewayCompleteLoginService,
			// Keep the original login service for backward compatibility if needed
			gateway.NewGatewayLoginService,
			// Other services
			gateway.NewGatewayLogoutService,
			// gateway.NewGatewaySendVerifyEmailService,
			gateway.NewGatewayRefreshTokenService,
			// gateway.NewGatewayResetPasswordService,
			// gateway.NewGatewayForgotPasswordService,
			// me.NewGetMeService,
			// me.NewUpdateMeService,
			// me.NewVerifyProfileService,
			// me.NewDeleteMeService,
		),
	)
}
