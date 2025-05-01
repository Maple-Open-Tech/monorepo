package mapleauth

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		http.Module(),
	)
}
