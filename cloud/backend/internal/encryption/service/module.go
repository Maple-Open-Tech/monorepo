// cloud/backend/internal/encryption/service/module.go
package service

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/service/encryptedfile"
)

// Module registers all services for encrypted files
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			encryptedfile.NewCreateEncryptedFileService,
			encryptedfile.NewGetEncryptedFileByIDService,
			encryptedfile.NewGetEncryptedFileByFileIDService,
			encryptedfile.NewUpdateEncryptedFileService,
			encryptedfile.NewDeleteEncryptedFileService,
			encryptedfile.NewListEncryptedFilesService,
			encryptedfile.NewDownloadEncryptedFileService,
			encryptedfile.NewGetEncryptedFileDownloadURLService,
		),
	)
}
