package maplesend

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/common"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
		),
	)
}
