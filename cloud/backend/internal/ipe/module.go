package ipe

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		http.Module(),
	)
}
