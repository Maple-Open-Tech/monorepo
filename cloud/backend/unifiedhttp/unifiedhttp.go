// github.com/comiccoin-network/monorepo/cloud/comiccoin/unifiedhttp/unifiedhttp.go
package unifiedhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func NewUnifiedHTTPServer(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{Addr: ":8000"}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
