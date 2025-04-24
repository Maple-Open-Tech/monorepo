package maplesend

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			repo.Module(),
			usecase.Module(),
			service.Module(),
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
			http.Module(),
		),
	)
}
