package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/repo/baseuser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/repo/templatedemailer"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewRepository,
			baseuser.NewRepository,
			templatedemailer.NewTemplatedEmailer,
		),
	)
}
