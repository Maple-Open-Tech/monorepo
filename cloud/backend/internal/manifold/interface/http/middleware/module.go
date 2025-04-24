// internal/manifold/interface/http/middleware/module.go
package middleware

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/blacklist"
	ipcb "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/ipcountryblocker"
)

// Module provides middleware components
func Module() fx.Option {
	return fx.Provide(
		blacklist.NewProvider,
		ipcb.NewProvider,
		NewMiddleware,
	)
}
