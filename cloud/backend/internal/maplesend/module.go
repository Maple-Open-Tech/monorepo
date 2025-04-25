package maplesend

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		http.Module(),
	)
}
