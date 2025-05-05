package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/repo/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/repo/templatedemailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/repo/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewRepository,
			user.NewRepository,
			templatedemailer.NewTemplatedEmailer,
			incomeproperty.NewMongoRepository,
		),
	)
}
