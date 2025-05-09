package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/service/me"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			me.NewDeleteMeService,
			me.NewGetMeService,
			me.NewUpdateMeService,
			me.NewVerifyProfileService,
		),
	)
}
