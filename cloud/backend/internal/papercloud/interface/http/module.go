package http

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/interface/http/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/interface/http/middleware"
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
			unifiedhttp.AsRoute(commonhttp.NewGetIncomePropertyEvaluatorVersionHTTPHandler),
		),
	)
}
