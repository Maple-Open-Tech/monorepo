package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/repo/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/repo/templatedemailer"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewRepository,
			federateduser.NewRepository,

			// Annotate the constructor to specify which parameter should receive the named dependency
			fx.Annotate(
				templatedemailer.NewTemplatedEmailer,
				fx.ParamTags(`name:"papercloud-module-emailer"`, `name:"maplesend-module-emailer"`),
			),
		),
	)
}
