package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/cmd/initialize"
	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/cmd/remote"
	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/cmd/version"
	// pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/common/preferences"
)

// var (
// 	preferences *pref.Preferences
// )

// // Initialize function will be called when every command gets called.
// func init() {
// 	preferences = pref.PreferencesInstance()
// }

var rootCmd = &cobra.Command{
	Use:   "papercloud-cli",
	Short: "PaperCloud CLI",
	Long:  `PaperCloud Command Line Interface`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(version.VersionCmd())
	rootCmd.AddCommand(initialize.InitializeCmd())
	rootCmd.AddCommand(remote.RemoteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
