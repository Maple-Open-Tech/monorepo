// internal/manifold/interface/http/module.go
package http

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/manifold/interface/http/middleware"
)

func Module() fx.Option {
	return fx.Options(
		middleware.Module(), // Include middleware module
		fx.Provide(
			NewUnifiedHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
		),
		fx.Provide(
			AsRoute(NewEchoHandler),
			AsRoute(NewGetHealthCheckHTTPHandler),
			// Add other routes here
		),
		fx.Invoke(func(*http.Server) {}),
	)
}
