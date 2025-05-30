package pkg

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/distributedmutex"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/emailer/mailgun"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/blacklist"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/ipcountryblocker"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodb"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/object/s3"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				mailgun.NewPaperCloudPropertyEvaluatorModuleEmailer,
				fx.ResultTags(`name:"papercloud-module-emailer"`), // Create name for better dependency management handling.
			),
			fx.Annotate(
				mailgun.NewPaperCloudPropertyEvaluatorModuleEmailer, //TODO: TEMPORARILY USED AS AN EXAMPLE.
				fx.ResultTags(`name:"maplesend-module-emailer"`),
			),
		),
		fx.Provide(
			blacklist.NewProvider,
			distributedmutex.NewAdapter,
			ipcountryblocker.NewProvider,
			jwt.NewProvider,
			password.NewProvider,
			mongodb.NewProvider,
			mongodbcache.NewProvider,
			s3.NewProvider,
		),
	)
}
