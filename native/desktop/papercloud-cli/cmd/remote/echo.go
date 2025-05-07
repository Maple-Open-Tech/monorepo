// native/desktop/papercloud-cli/cmd/remote/echo.go
package remote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

func EchoCmd() *cobra.Command {
	var text string

	var cmd = &cobra.Command{
		Use:   "echo",
		Short: "Echo text to backend",
		Long:  `Command will execute submitting any text to the server and the server will respond back with the text you submitted`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echoing text to server...")

			// If no text provided via flag, get from args
			if text == "" && len(args) > 0 {
				text = strings.Join(args, " ")
			}

			// If still no text, prompt the user
			if text == "" {
				fmt.Println("No text provided. Please provide text to echo using --text flag or as arguments.")
				return
			}

			preferences := pref.PreferencesInstance()
			serverURL := preferences.CloudProviderAddress
			if serverURL == "" {
				serverURL = "http://localhost:8000" // Default if not configured
			}

			// Make a POST request to the echo endpoint
			echoURL := fmt.Sprintf("%s/echo", serverURL)
			fmt.Printf("Connecting to: %s\n", echoURL)

			resp, err := http.Post(echoURL, "text/plain", bytes.NewBuffer([]byte(text)))
			if err != nil {
				fmt.Printf("Error connecting to server: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Check if the response was successful
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Server returned error status: %s\n", resp.Status)
				return
			}

			// Read and display the response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			// Display the echoed text
			fmt.Printf("Server echoed: %s\n", string(body))
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&text, "text", "t", "", "Text to echo")

	return cmd
}
