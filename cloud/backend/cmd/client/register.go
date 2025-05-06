package client

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/e2ee"
)

// RegisterUserCmd creates a new Cobra command for registering a user account.
// This command requires email, password, first name, and last name as mandatory flags.
// Optional flags include timezone, country, phone number, beta access code, agreement flags, and module.
//
// Usage:
//
//	./backend client register --email <email> --password <password> --firstname <firstname> --lastname <lastname> [flags]
//
// Example (Required flags only):
//
//	./backend client register -e user@example.com -p securepass -f John -l Doe
//
// Example (With optional flags):
//
//	./backend client register --email user@example.com --password securepass --firstname Jane --lastname Smith --timezone "America/New_York" --country USA --phone "123-456-7890" --beta-code BETA123 --agree-terms --agree-promotions --agree-tracking --module 1
//
// Note: If timezone is not provided, it defaults to "UTC". If country is not provided, it defaults to "Canada". Agreement flags default to false. Module defaults to 0.
func RegisterUserCmd() *cobra.Command {
	var email, password, firstName, lastName, timezone, country, phone, betaAccessCode string
	var agreeTerms, agreePromotions, agreeTracking bool
	var module int

	var cmd = &cobra.Command{
		Use:   "register",
		Short: "Register user account",
		Long: `Register a new user account in the system.

This command requires you to provide an email, password, first name, and last name.
You can optionally provide timezone, country, phone number, a beta access code,
specify agreement to terms, promotions, and tracking, and specify the registration module.

Examples:
		# Register with only required fields
		register --email user@example.com --password mysecret --firstname John --lastname Doe

		# Register with all fields using short flags (note: only some have short flags)
		register -e test@domain.com -p pass123 -f Jane -l Smith -t "America/Toronto" -c "USA" -n "555-1234" --beta-code ABCDE --agree-terms --module 2

		# Register using a mix of short and long flags, enabling all agreements and specifying module
		register --email another@user.net -p anotherpass -f Bob -l Williams --timezone "Europe/London" --agree-terms --agree-promotions --agree-tracking --module 1`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Registering...")

			// Required flags are enforced by Cobra.
			// Default values for timezone, country, and module are handled by flag definitions.

			// Call registerUser with all collected flags
			err := registerUser(email, password, firstName, lastName, timezone, country, phone, betaAccessCode, agreeTerms, agreePromotions, agreeTracking, module)
			if err != nil {
				log.Fatalf("Failed to register user: %v", err)
			}
			log.Println("User registration process initiated.") // Assuming registerUser is async or just starts the process
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address for the user (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user (required)")
	cmd.Flags().StringVarP(&firstName, "firstname", "f", "", "First name for the user (required)")
	cmd.Flags().StringVarP(&lastName, "lastname", "l", "", "Last name for the user (required)")
	cmd.Flags().StringVarP(&timezone, "timezone", "t", "UTC", "Timezone for the user (e.g., \"America/New_York\")")
	cmd.Flags().StringVarP(&country, "country", "c", "Canada", "Country for the user")
	cmd.Flags().StringVarP(&phone, "phone", "n", "", "Phone number for the user")
	cmd.Flags().StringVar(&betaAccessCode, "beta-code", "", "Beta access code (if required)")
	cmd.Flags().BoolVar(&agreeTerms, "agree-terms", false, "Agree to Terms of Service")
	cmd.Flags().BoolVar(&agreePromotions, "agree-promotions", false, "Agree to receive promotions")
	cmd.Flags().BoolVar(&agreeTracking, "agree-tracking", false, "Agree to tracking across third-party apps and services")
	cmd.Flags().IntVarP(&module, "module", "m", 0, "Module the user is registering for") // Added module flag

	// Mark required flags - Cobra will handle enforcing these
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("firstname")
	cmd.MarkFlagRequired("lastname")

	return cmd
}

// registerUser handles the actual user registration logic.
func registerUser(email, plainPassword, firstName, lastName, timezone, country, phone, betaAccessCode string, agreeTerms, agreePromotions, agreeTracking bool, module int) error { // Added module parameter
	cfg := config.NewProvider()

	// Create a new E2EE client
	client := e2ee.NewClient(e2ee.ClientConfig{
		ServerURL: fmt.Sprintf("%s://%s:%s", "http", cfg.App.IP, cfg.App.Port),
	})

	// Log the received parameters (for debugging/demonstration)
	log.Printf("Attempting registration with:")
	log.Printf("  Email: %s", email)
	log.Printf("  Password: [REDACTED]")
	log.Printf("  First Name: %s", firstName)
	log.Printf("  Last Name: %s", lastName)
	log.Printf("  Timezone: %s (Note: Not sent in current Register call)", timezone)
	log.Printf("  Country: %s", country)
	log.Printf("  Phone: %s", phone)
	log.Printf("  Beta Code: %s", betaAccessCode)
	log.Printf("  Agree Terms: %t", agreeTerms)
	log.Printf("  Agree Promotions: %t", agreePromotions)
	log.Printf("  Agree Tracking: %t", agreeTracking)
	log.Printf("  Module: %d", module) // Log the module

	// Call the Register method with all required parameters according to its signature.
	recoveryKey, err := client.Register(
		email, plainPassword,
		betaAccessCode, firstName, lastName, phone, country, timezone,
		agreeTerms, agreePromotions, agreeTracking,
		module, // Pass the module value
	)
	if err != nil {
		// Log the specific error from the registration call
		log.Printf("e2ee.Client.Register failed: %v", err)
		// Return a more general error to the caller in Run
		return fmt.Errorf("registration via E2EE client failed: %w", err)
	}

	// Display recovery key information
	recoveryKeyInfo := client.GetRecoveryKeyInfo(recoveryKey)
	fmt.Println("Registration successful. Recovery Key Info:")
	fmt.Println(recoveryKeyInfo)

	return nil
}
