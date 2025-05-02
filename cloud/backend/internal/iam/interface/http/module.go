package http

import (
	"go.uber.org/fx"

	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/interface/http/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/interface/http/middleware"
	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			middleware.NewMiddleware,
		),
		fx.Provide(
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayUserRegisterHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayLoginHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayLogoutHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayRefreshTokenHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayResetPasswordHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayForgotPasswordHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayVerifyEmailHTTPHandler),
		),
	)
}
