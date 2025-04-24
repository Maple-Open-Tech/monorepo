package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service/gateway"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			gateway.NewGatewayUserRegisterService,
		),
	)
}
