// cloud/backend/internal/vault/repo/module.go
package repo

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/repo/encryptedfile"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			encryptedfile.NewRepository,
		),
	)
}
