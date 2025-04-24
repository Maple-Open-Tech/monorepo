package http

import (
	"net/http"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
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
		),
		fx.Invoke(func(*http.Server) {}),
	)
}
