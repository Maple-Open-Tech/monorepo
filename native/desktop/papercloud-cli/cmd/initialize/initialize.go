package initialize

import (
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

// Command line argument flags
var (
	flagDataDirectory        string
	flagCloudProviderAddress string
)

func InitializeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize your local PaperCloud for the first time. Note: This action cannot be undone once executed.",
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewDevelopment()

			if preferences.DataDirectory != "" {
				log.Fatalf("You have already configured PaperCloud: DataDirectory was set with: %v\n", preferences.DataDirectory)
			}
			preferences.SetDataDirectory(flagDataDirectory)

			if preferences.CloudProviderAddress != "" {
				log.Fatalf("You have already configured PaperCloud: CloudProviderAddress was set with: %v\n", preferences.CloudProviderAddress)
			}
			preferences.SetCloudProviderAddress(flagCloudProviderAddress)

			logger.Debug("Configued PaperCloud",
				zap.Any("DataDirectory", preferences.DataDirectory),
				zap.Any("CloudProviderAddress", preferences.CloudProviderAddress),
				zap.Any("FilePathPreferences", preferences.GetFilePathOfPreferencesFile()))
		},
	}
	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", pref.GetDefaultDataDirectory(), "The data directory to save to")
	cmd.Flags().StringVar(&flagCloudProviderAddress, "cloud-provider-address", pref.GetDefaultCloudProviderAddress(), "The address of the cloud provider for PaperCloud")

	return cmd
}
