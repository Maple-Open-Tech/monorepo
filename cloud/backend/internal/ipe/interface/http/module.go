package http

import (
	"go.uber.org/fx"

	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/middleware"
	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
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
