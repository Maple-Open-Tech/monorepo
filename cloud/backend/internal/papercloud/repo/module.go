package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/templatedemailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/user"
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
