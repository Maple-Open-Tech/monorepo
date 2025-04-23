// github.com/Maple-Open-Tech/monorepo/cloud/backend/cmd/daemon/daemon.go
package daemon

import (
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

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
		fx.Provide(unifiedhttp.NewUnifiedHTTPServer),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
