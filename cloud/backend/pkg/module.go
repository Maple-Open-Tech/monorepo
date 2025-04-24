package pkg

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/distributedmutex"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/blacklist"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/ipcountryblocker"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodb"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			distributedmutex.NewAdapter,
			blacklist.NewProvider,
			ipcountryblocker.NewProvider,
			jwt.NewProvider,
			password.NewProvider,
			mongodb.NewProvider,
			mongodb.NewProvider,
			mongodbcache.NewProvider,
		),
	)
}
