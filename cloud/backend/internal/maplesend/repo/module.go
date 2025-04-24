package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo/templatedemailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewRepository,
			user.NewRepository,
			templatedemailer.NewTemplatedEmailer,
		),
	)
}
