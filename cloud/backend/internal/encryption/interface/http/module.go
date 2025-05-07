// cloud/backend/internal/encryption/interface/http/module.go
package http

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/interface/http/encryptedfile"
	unifiedhttp "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http"
)

// Module registers all HTTP handlers for encrypted files
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			unifiedhttp.AsRoute(encryptedfile.NewCreateEncryptedFileHandler),
			unifiedhttp.AsRoute(encryptedfile.NewGetEncryptedFileByIDHandler),
			unifiedhttp.AsRoute(encryptedfile.NewGetEncryptedFileByFileIDHandler),
			// // unifiedhttp.AsRoute(encryptedfile.NewUpdateEncryptedFileHandler),
			unifiedhttp.AsRoute(encryptedfile.NewDeleteEncryptedFileHandler),
			unifiedhttp.AsRoute(encryptedfile.NewListEncryptedFilesHandler),
			// unifiedhttp.AsRoute(encryptedfile.NewDownloadEncryptedFileHandler),
			// unifiedhttp.AsRoute(encryptedfile.NewGetEncryptedFileDownloadURLHandler),
		),
	)
}
