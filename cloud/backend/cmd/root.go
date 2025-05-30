// github.com/Maple-Open-Tech/monorepo/cloud/backend/cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/cmd/daemon"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/cmd/version"
)

// Initialize function will be called when every command gets called.
func init() {

}

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "MOT Cloud Backend Services",
	Long:  `Maple Open Technologies Cloud Backend Services`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(daemon.DaemonCmd())
	rootCmd.AddCommand(version.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
