// cloud/backend/internal/encryption/module.go
package encryption

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		// Will add this in the future:
		// http.Module(),
	)
}
