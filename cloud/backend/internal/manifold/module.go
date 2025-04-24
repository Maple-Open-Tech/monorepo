package manifold

import (
	"net/http"

	"go.uber.org/fx"

	commonhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg"
)

func Module() fx.Option {
	return fx.Options(
		pkg.Module(),
		maplesend.Module(),
		commonhttp.Module(),
		fx.Invoke(func(*http.Server) {}),
	)
}
