package http

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http/middleware"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			middleware.NewMiddleware,
		),
		fx.Provide(
			unifiedhttp.AsRoute(me.NewGetMeHTTPHandler),
			unifiedhttp.AsRoute(me.NewPutUpdateMeHTTPHandler),
			unifiedhttp.AsRoute(me.NewDeleteMeHTTPHandler),
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayUserRegisterHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayLoginHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayLogoutHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayRefreshTokenHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayResetPasswordHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayForgotPasswordHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayVerifyEmailHTTPHandler),
			unifiedhttp.AsRoute(me.NewPostVerifyProfileHTTPHandler),
		),
	)
}
