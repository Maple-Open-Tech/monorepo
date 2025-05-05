// cloud/backend/internal/ipe/service/module.go
package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/service/me"
	// "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/service/person/owner"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			// Gateway services (existing)
			gateway.NewGatewayUserRegisterService,
			gateway.NewGatewayLoginService,
			gateway.NewGatewayLogoutService,
			gateway.NewGatewayResetPasswordService,
			gateway.NewGatewaySendVerifyEmailService,
			gateway.NewGatewayVerifyEmailService,
			gateway.NewGatewayRefreshTokenService,
			gateway.NewGatewayForgotPasswordService,

			// Me services (existing)
			me.NewGetMeService,
			me.NewUpdateMeService,
			me.NewVerifyProfileService,
			me.NewDeleteMeService,
		),
	)
}
