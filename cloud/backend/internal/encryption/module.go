// cloud/backend/internal/encryption/module.go
package encryption

import (
	"go.uber.org/fx"

	iface "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/interface/http"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/repo"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/service"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase"
)

func Module() fx.Option {
	return fx.Options(
		repo.Module(),
		usecase.Module(),
		service.Module(),
		iface.Module(),
	)
}
