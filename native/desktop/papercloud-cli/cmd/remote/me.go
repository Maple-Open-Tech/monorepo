// native/desktop/papercloud-cli/cmd/remote/me.go
package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// MeCmd returns a cobra command that retrieves the user's profile
func MeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "me",
		Short: "Get user profile information",
		Long: `
Retrieves and displays the current user's profile information from the server.
This command requires you to be logged in with a valid access token.

Example:
  papercloud-cli remote me
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Fetching user profile...")

			// Get preferences with stored tokens
			preferences := pref.PreferencesInstance()

			// Check if user is logged in
			if preferences.LoginResponse == nil || preferences.LoginResponse.AccessToken == "" {
				fmt.Println("Error: You are not logged in. Please login first.")
				return
			}

			// Check if token is expired
			if time.Now().After(preferences.LoginResponse.AccessTokenExpiryTime) {
				fmt.Println("Error: Your access token has expired. Please login again.")
				return
			}

			// Get server URL
			serverURL := preferences.CloudProviderAddress
			if serverURL == "" {
				serverURL = "http://localhost:8000" // Default if not configured
				fmt.Println("Warning: Cloud provider address not configured. Using default:", serverURL)
			}

			// Create request
			meURL := fmt.Sprintf("%s/papercloud/api/v1/me", serverURL)
			req, err := http.NewRequest("GET", meURL, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}

			// Add authorization header with JWT token format
			// The backend expects "JWT <token>" format based on the middleware code
			authHeader := fmt.Sprintf("JWT %s", preferences.LoginResponse.AccessToken)
			req.Header.Set("Authorization", authHeader)

			// Execute request
			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error connecting to server: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Check response status
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Server returned error status: %s\n", resp.Status)

				// Try to read error message
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("Error details: %s\n", string(body))

				// Additional help message for common errors
				if resp.StatusCode == http.StatusUnauthorized {
					fmt.Println("\nYour authentication token may be invalid or expired. Please try logging in again.")
				}
				return
			}

			// Read and parse response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			// Define response structure based on MeResponseDTO in the backend
			type MeResponse struct {
				ID               string    `json:"id"`
				Email            string    `json:"email"`
				FirstName        string    `json:"first_name"`
				LastName         string    `json:"last_name"`
				Name             string    `json:"name"`
				LexicalName      string    `json:"lexical_name"`
				Role             int       `json:"role"`
				WasEmailVerified bool      `json:"was_email_verified"`
				Phone            string    `json:"phone,omitempty"`
				Country          string    `json:"country,omitempty"`
				Timezone         string    `json:"timezone"`
				Region           string    `json:"region,omitempty"`
				City             string    `json:"city,omitempty"`
				PostalCode       string    `json:"postal_code,omitempty"`
				AddressLine1     string    `json:"address_line1,omitempty"`
				AddressLine2     string    `json:"address_line2,omitempty"`
				AgreePromotions  bool      `json:"agree_promotions,omitempty"`
				AgreeToTracking  bool      `json:"agree_to_tracking_across_third_party_apps_and_services,omitempty"`
				CreatedAt        time.Time `json:"created_at,omitempty"`
				Status           int       `json:"status"`
			}

			var profile MeResponse
			if err := json.Unmarshal(body, &profile); err != nil {
				fmt.Printf("Error parsing profile data: %v\n", err)
				fmt.Printf("Raw response: %s\n", string(body))
				return
			}

			// Display profile information in a formatted way
			fmt.Println("\n=== User Profile ===")
			fmt.Printf("ID: %s\n", profile.ID)
			fmt.Printf("Name: %s %s\n", profile.FirstName, profile.LastName)
			fmt.Printf("Email: %s\n", profile.Email)

			// Map role to human-readable string
			roleMap := map[int]string{
				1: "Root",
				2: "Company",
				3: "Individual",
			}
			roleStr := roleMap[profile.Role]
			if roleStr == "" {
				roleStr = fmt.Sprintf("Unknown (%d)", profile.Role)
			}
			fmt.Printf("Role: %s\n", roleStr)

			// Map status to human-readable string
			statusMap := map[int]string{
				1:   "Active",
				50:  "Locked",
				100: "Archived",
			}
			statusStr := statusMap[profile.Status]
			if statusStr == "" {
				statusStr = fmt.Sprintf("Unknown (%d)", profile.Status)
			}
			fmt.Printf("Status: %s\n", statusStr)

			// Display contact information
			fmt.Println("\n--- Contact Information ---")
			fmt.Printf("Phone: %s\n", profile.Phone)
			fmt.Printf("Country: %s\n", profile.Country)
			if profile.Region != "" {
				fmt.Printf("Region: %s\n", profile.Region)
			}
			if profile.City != "" {
				fmt.Printf("City: %s\n", profile.City)
			}
			if profile.AddressLine1 != "" {
				fmt.Printf("Address: %s\n", profile.AddressLine1)
				if profile.AddressLine2 != "" {
					fmt.Printf("         %s\n", profile.AddressLine2)
				}
			}
			if profile.PostalCode != "" {
				fmt.Printf("Postal Code: %s\n", profile.PostalCode)
			}
			fmt.Printf("Timezone: %s\n", profile.Timezone)

			// Display preferences
			fmt.Println("\n--- Preferences ---")
			fmt.Printf("Receives Promotions: %t\n", profile.AgreePromotions)
			fmt.Printf("Agrees to Tracking: %t\n", profile.AgreeToTracking)

			// Display account info
			fmt.Println("\n--- Account Information ---")
			fmt.Printf("Email Verified: %t\n", profile.WasEmailVerified)
			if !profile.CreatedAt.IsZero() {
				fmt.Printf("Account Created: %s\n", profile.CreatedAt.Format("January 2, 2006"))
			}
		},
	}

	return cmd
}
