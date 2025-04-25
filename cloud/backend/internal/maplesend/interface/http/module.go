package http

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			middleware.NewMiddleware,
		),
		fx.Provide(
			// These are needed as dependencies but not as routes
			me.NewGetMeHTTPHandler,
			me.NewPutUpdateMeHTTPHandler,
			me.NewDeleteMeHTTPHandler,

			// This is the combined handler that will be registered as a route
			unifiedhttp.AsRoute(me.NewMeCombinedHandler),

			// Other routes
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
