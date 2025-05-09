// cloud/backend/internal/papercloud/repo/module.go
package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/collection"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/file"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/templatedemailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bannedipaddress.NewRepository,
			collection.NewRepository,
			file.NewFileMetadataRepository,
			file.NewFileStorageRepository,
			file.NewFileRepository,
			user.NewRepository,
			templatedemailer.NewTemplatedEmailer,
		),
	)
}
