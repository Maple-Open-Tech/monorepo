package iam

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		http.Module(),
	)
}
