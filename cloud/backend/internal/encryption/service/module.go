// cloud/backend/internal/encryption/service.go
package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/service/encryptedfile"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			encryptedfile.NewEncryptedFileService,
		),
	)
}
