package manifold

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe"
	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg"
)

func Module() fx.Option {
	return fx.Options(
		pkg.Module(),
		commonhttp.Module(),
		iam.Module(),
		ipe.Module(),
		maplesend.Module(),
		fx.Invoke(func(*http.Server) {}),
	)
}
