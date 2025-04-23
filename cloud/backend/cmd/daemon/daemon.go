// github.com/Maple-Open-Tech/monorepo/cloud/backend/cmd/daemon/daemon.go
package daemon

import (
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/unifiedhttp"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run the cloud-services backend",
		Run: func(cmd *cobra.Command, args []string) {
			doRunDaemon()
		},
	}
	return cmd
}

func doRunDaemon() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			unifiedhttp.NewUnifiedHTTPServer,
			fx.Annotate(
				unifiedhttp.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
		),
		fx.Provide(
			unifiedhttp.AsRoute(unifiedhttp.NewEchoHandler),
			unifiedhttp.AsRoute(unifiedhttp.NewGetHealthCheckHTTPHandler),
			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
