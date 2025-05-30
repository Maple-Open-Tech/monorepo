// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser/model.go
package federateduser

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FederatedUserStatusActive   = 1   // User is active and can log in.
	FederatedUserStatusLocked   = 50  // User account is locked, typically due to too many failed login attempts.
	FederatedUserStatusArchived = 100 // User account is archived and cannot log in.

	FederatedUserRoleRoot       = 1 // Root user, has all permissions
	FederatedUserRoleCompany    = 2 // Company user, has permissions for company-related operations
	FederatedUserRoleIndividual = 3 // Individual user, has permissions for individual-related operations
)

type FederatedUser struct {
	ID                                             primitive.ObjectID `bson:"_id" json:"id"`
	Email                                          string             `bson:"email" json:"email"`
	FirstName                                      string             `bson:"first_name" json:"first_name"`
	LastName                                       string             `bson:"last_name" json:"last_name"`
	Name                                           string             `bson:"name" json:"name"`
	LexicalName                                    string             `bson:"lexical_name" json:"lexical_name"`
	Role                                           int8               `bson:"role" json:"role"`
	Status                                         int8               `bson:"status" json:"status"`
	WasEmailVerified                               bool               `bson:"was_email_verified" json:"was_email_verified,omitempty"`
	EmailVerificationCode                          string             `bson:"email_verification_code,omitempty" json:"email_verification_code,omitempty"`
	EmailVerificationExpiry                        time.Time          `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry"`
	PasswordResetVerificationCode                  string             `bson:"password_reset_verification_code,omitempty" json:"password_reset_verification_code,omitempty"`
	PasswordResetVerificationExpiry                time.Time          `bson:"password_reset_verification_expiry,omitempty" json:"password_reset_verification_expiry"`
	Phone                                          string             `bson:"phone" json:"phone,omitempty"`
	Country                                        string             `bson:"country" json:"country,omitempty"`
	Timezone                                       string             `bson:"timezone" json:"timezone"`
	Region                                         string             `bson:"region" json:"region,omitempty"`
	City                                           string             `bson:"city" json:"city,omitempty"`
	PostalCode                                     string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1                                   string             `bson:"address_line1" json:"address_line1,omitempty"`
	AddressLine2                                   string             `bson:"address_line2" json:"address_line2,omitempty"`
	HasShippingAddress                             bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName                                   string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone                                  string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry                                string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion                                 string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity                                   string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode                             string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1                           string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2                           string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	AgreeTermsOfService                            bool               `bson:"agree_terms_of_service" json:"agree_terms_of_service,omitempty"`
	AgreePromotions                                bool               `bson:"agree_promotions" json:"agree_promotions,omitempty"`
	AgreeToTrackingAcrossThirdPartyAppsAndServices bool               `bson:"agree_to_tracking_across_third_party_apps_and_services" json:"agree_to_tracking_across_third_party_apps_and_services,omitempty"`

	// --- E2EE Related ---
	Salt                              string `json:"salt"`
	PublicKey                         string `json:"publicKey"`
	EncryptedMasterKey                string `json:"encryptedMasterKey"`
	EncryptedPrivateKey               string `json:"encryptedPrivateKey"`
	EncryptedRecoveryKey              string `json:"encryptedRecoveryKey"`
	MasterKeyEncryptedWithRecoveryKey string `json:"masterKeyEncryptedWithRecoveryKey"`
	VerificationID                    string `json:"verificationID"`

	// --- Metadata ---
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address"`
	CreatedByUserID       primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	CreatedByName         string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address"`
	ModifiedByUserID      primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedAt            time.Time          `bson:"modified_at" json:"modified_at"`
	ModifiedByName        string             `bson:"modified_by_name" json:"modified_by_name"`

	// OTPEnabled controls whether we force 2FA or not during login.
	OTPEnabled bool `bson:"otp_enabled" json:"otp_enabled"`

	// OTPVerified indicates user has successfully validated their opt token afer enabling 2FA thus turning it on.
	OTPVerified bool `bson:"otp_verified" json:"otp_verified"`

	// OTPValidated automatically gets set as `false` on successful login and then sets `true` once successfully validated by 2FA.
	OTPValidated bool `bson:"otp_validated" json:"otp_validated"`

	// OTPSecret the unique one-time password secret to be shared between our
	// backend and 2FA authenticator sort of apps that support `TOPT`.
	OTPSecret string `bson:"otp_secret" json:"-"`

	// OTPAuthURL is the URL used to share.
	OTPAuthURL string `bson:"otp_auth_url" json:"-"`

	// OTPBackupCodeHash is the one-time use backup code which resets the 2FA settings and allow the user to setup 2FA from scratch for the user.
	OTPBackupCodeHash string `bson:"otp_backup_code_hash" json:"-"`

	// OTPBackupCodeHashAlgorithm tracks the hashing algorithm used.
	OTPBackupCodeHashAlgorithm string `bson:"otp_backup_code_hash_algorithm" json:"-"`
}

// FederatedUserFilter represents the filter criteria for listing users
type FederatedUserFilter struct {
	// Basic filters
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Role   int8    `json:"role,omitempty"`
	Status int8    `json:"status,omitempty"`

	// Date range filters
	CreatedAtStart *time.Time `json:"created_at_start,omitempty"`
	CreatedAtEnd   *time.Time `json:"created_at_end,omitempty"`

	// Pagination - cursor based
	LastID        *primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt *time.Time          `json:"last_created_at,omitempty"`
	Limit         int64               `json:"limit,omitempty"`

	// Search term for text search across multiple fields
	SearchTerm *string `json:"search_term,omitempty"`
}

// FederatedUserFilterResult represents the result of a filtered list operation
type FederatedUserFilterResult struct {
	Users         []*FederatedUser   `json:"users"`
	HasMore       bool               `json:"has_more"`
	LastID        primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt time.Time          `json:"last_created_at"`
	TotalCount    uint64             `json:"total_count,omitempty"`
}
