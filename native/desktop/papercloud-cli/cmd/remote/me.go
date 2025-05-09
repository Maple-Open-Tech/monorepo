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

// MeResponse defines the structure of the me endpoint response
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

// executeGetMeRequest performs the HTTP request to get user profile
// Returns the profile data and a boolean indicating if a retry should be attempted
func executeGetMeRequest(retryAfterRefresh bool) (*MeResponse, bool, error) {
	preferences := pref.PreferencesInstance()

	// Check if user is logged in
	if preferences.LoginResponse == nil || preferences.LoginResponse.AccessToken == "" {
		return nil, false, fmt.Errorf("you are not logged in. Please login first")
	}

	// Check if token is expired - if so, try refresh unless we're already retrying
	if time.Now().After(preferences.LoginResponse.AccessTokenExpiryTime) {
		if retryAfterRefresh {
			return nil, false, fmt.Errorf("access token expired and refresh attempt already failed")
		}

		fmt.Println("Access token expired. Attempting to refresh...")
		success, err := RefreshTokens()
		if err != nil && success {
			return nil, true, nil // Signal to retry with new token
		}
		return nil, false, fmt.Errorf("access token expired and refresh failed: %v", err)
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
		return nil, false, fmt.Errorf("error creating request: %v", err)
	}

	// Add authorization header with JWT token
	authHeader := fmt.Sprintf("JWT %s", preferences.LoginResponse.AccessToken)
	req.Header.Set("Authorization", authHeader)

	// Execute request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("error connecting to server: %v", err)
	}
	defer resp.Body.Close()

	// Handle 401 Unauthorized - could be expired token
	if resp.StatusCode == http.StatusUnauthorized {
		if retryAfterRefresh {
			return nil, false, fmt.Errorf("authentication failed even after token refresh")
		}

		fmt.Println("Authentication failed. Attempting to refresh token...")
		success, err := RefreshTokens()
		if err != nil || !success {
			return nil, false, fmt.Errorf("token refresh failed: %v", err)
		}
		return nil, true, nil // Signal to retry with new token
	}

	// Check other error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("server returned error status: %s - %s",
			resp.Status, string(body))
	}

	// Read and parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response: %v", err)
	}

	var profile MeResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, false, fmt.Errorf("error parsing profile data: %v\nRaw response: %s",
			err, string(body))
	}

	return &profile, false, nil
}

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

			// First attempt
			profile, shouldRetry, err := executeGetMeRequest(false)

			// If we need to retry after token refresh
			if shouldRetry {
				fmt.Println("Retrying with new access token...")
				profile, _, err = executeGetMeRequest(true)
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Display profile information
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
