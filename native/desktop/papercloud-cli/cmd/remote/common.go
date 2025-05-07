// native/desktop/papercloud-cli/cmd/remote/common.go
package remote

import (
	"fmt"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/e2ee"
	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// createE2EEClient creates a new e2ee client with server URL from preferences
func createE2EEClient() *e2ee.Client {
	preferences := pref.PreferencesInstance()
	serverURL := preferences.CloudProviderAddress
	if serverURL == "" {
		serverURL = "http://localhost:8000" // Default if not configured
		fmt.Println("Warning: Cloud provider address not configured. Using default:", serverURL)
	}

	return e2ee.NewClient(e2ee.ClientConfig{
		ServerURL: serverURL,
	})
}
