package e2ee

// RegistrationPayload contains all data sent to the server during registration
// All sensitive data is encrypted before being added to this struct
type RegistrationPayload struct {
	// --- Application and personal identiable information (PII) related ---
	BetaAccessCode                                 string `json:"beta_access_code"` // Temporary code for beta access
	FirstName                                      string `json:"first_name"`
	LastName                                       string `json:"last_name"`
	Email                                          string `json:"email"`
	Phone                                          string `json:"phone,omitempty"`
	Country                                        string `json:"country,omitempty"`
	Timezone                                       string `json:"timezone"`
	AgreeTermsOfService                            bool   `json:"agree_terms_of_service,omitempty"`
	AgreePromotions                                bool   `json:"agree_promotions,omitempty"`
	AgreeToTrackingAcrossThirdPartyAppsAndServices bool   `json:"agree_to_tracking_across_third_party_apps_and_services,omitempty"`

	// Module refers to which module the user is registering for.
	Module int `json:"module,omitempty"`

	// --- E2EE Related ---
	Salt                              string `json:"salt"`
	PublicKey                         string `json:"publicKey"`
	EncryptedMasterKey                string `json:"encryptedMasterKey"`
	EncryptedPrivateKey               string `json:"encryptedPrivateKey"`
	EncryptedRecoveryKey              string `json:"encryptedRecoveryKey"`
	MasterKeyEncryptedWithRecoveryKey string `json:"masterKeyEncryptedWithRecoveryKey"`
	VerificationID                    string `json:"verificationID"`
}

// KeySet represents the complete set of cryptographic keys for a user
// These are the unencrypted, sensitive keys that should never leave the client
type KeySet struct {
	MasterKey   []byte // Root key for the encryption system
	PublicKey   []byte // Public key (can be shared)
	PrivateKey  []byte // Private key (must be kept secret)
	RecoveryKey []byte // Backup key for account recovery
}

// ClientConfig holds the configuration for the E2EE client
type ClientConfig struct {
	// Base URL for the API server
	ServerURL string

	// Optional: Custom HTTP client
	HTTPClient HTTPClient
}

// Default configuration values
const (
	DefaultServerURL = "https://api.yourservice.com"
)
