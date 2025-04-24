package maplesend

import (
	"go.uber.org/fx"

	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/common"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		fx.Provide(
			middleware.NewMiddleware,
		),
		fx.Provide(
			unifiedhttp.AsRoute(commonhttp.NewGetMapleSendVersionHTTPHandler),
			unifiedhttp.AsRoute(gateway.NewGatewayUserRegisterHTTPHandler),
		),
	)
}
