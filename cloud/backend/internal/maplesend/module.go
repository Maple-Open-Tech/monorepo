package maplesend

import (
	"go.uber.org/fx"

	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/common"
	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/unifiedhttp"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
		),
	)
}
