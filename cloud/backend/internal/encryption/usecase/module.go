// cloud/backend/internal/encryption/usecase/module.go
package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
)

// Module registers all encrypted file use cases
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			encryptedfile.NewCreateEncryptedFileUseCase,
			encryptedfile.NewGetEncryptedFileByIDUseCase,
			encryptedfile.NewGetEncryptedFileByFileIDUseCase,
			encryptedfile.NewUpdateEncryptedFileUseCase,
			encryptedfile.NewDeleteEncryptedFileUseCase,
			encryptedfile.NewListEncryptedFilesUseCase,
			encryptedfile.NewDownloadEncryptedFileUseCase,
			encryptedfile.NewGetEncryptedFileDownloadURLUseCase,
		),
	)
}
