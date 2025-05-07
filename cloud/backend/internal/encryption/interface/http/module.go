// cloud/backend/internal/encryption/interface/http/module.go
package http

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http/encryptedfile"
)

// Module registers all HTTP handlers for encrypted files
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			encryptedfile.AsRoute(encryptedfile.NewCreateEncryptedFileHandler),
			encryptedfile.AsRoute(encryptedfile.NewGetEncryptedFileByIDHandler),
			encryptedfile.AsRoute(encryptedfile.NewGetEncryptedFileByFileIDHandler),
			encryptedfile.AsRoute(encryptedfile.NewUpdateEncryptedFileHandler),
			encryptedfile.AsRoute(encryptedfile.NewDeleteEncryptedFileHandler),
			encryptedfile.AsRoute(encryptedfile.NewListEncryptedFilesHandler),
			encryptedfile.AsRoute(encryptedfile.NewDownloadEncryptedFileHandler),
			encryptedfile.AsRoute(encryptedfile.NewGetEncryptedFileDownloadURLHandler),
		),
	)
}
