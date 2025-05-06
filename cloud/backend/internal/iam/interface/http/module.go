// cloud/backend/internal/iam/interface/http/module.go
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
			unifiedhttp.AsRoute(gateway.NewGatewayFederatedUserRegisterHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayVerifyEmailHTTPHandler),
			// Add the new E2EE login handlers
			unifiedhttp.AsRoute(gateway.NewGatewayRequestLoginOTTHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayVerifyLoginOTTHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayCompleteLoginHTTPHandler),
			// Keep the original login handler for backward compatibility if needed
			unifiedhttp.AsRoute(gateway.NewGatewayLoginHTTPHandler),
			// Other handlers
			unifiedhttp.AsRoute(gateway.NewGatewayLogoutHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayRefreshTokenHTTPHandler),
			// unifiedhttp.AsRoute(gateway.NewGatewayResetPasswordHTTPHandler),
			// unifiedhttp.AsRoute(gateway.NewGatewayForgotPasswordHTTPHandler),
		),
	)
}
