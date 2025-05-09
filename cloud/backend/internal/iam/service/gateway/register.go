package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
	uc_emailer "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/emailer"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/random"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayFederatedUserRegisterService interface {
	Execute(
		sessCtx context.Context,
		req *RegisterCustomerRequestIDO,
	) error
}

type gatewayFederatedUserRegisterServiceImpl struct {
	config                                    *config.Configuration
	logger                                    *zap.Logger
	passwordProvider                          password.Provider
	cache                                     mongodbcache.Cacher
	jwtProvider                               jwt.Provider
	userGetByEmailUseCase                     uc_user.FederatedUserGetByEmailUseCase
	userCreateUseCase                         uc_user.FederatedUserCreateUseCase
	userUpdateUseCase                         uc_user.FederatedUserUpdateUseCase
	sendFederatedUserVerificationEmailUseCase uc_emailer.SendFederatedUserVerificationEmailUseCase
}

func NewGatewayFederatedUserRegisterService(
	cfg *config.Configuration,
	logger *zap.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.FederatedUserGetByEmailUseCase,
	uc2 uc_user.FederatedUserCreateUseCase,
	uc3 uc_user.FederatedUserUpdateUseCase,
	uc4 uc_emailer.SendFederatedUserVerificationEmailUseCase,
) GatewayFederatedUserRegisterService {
	return &gatewayFederatedUserRegisterServiceImpl{cfg, logger, pp, cach, jwtp, uc1, uc2, uc3, uc4}
}

type RegisterCustomerRequestIDO struct {
	// --- Application and personal identiable information (PII) ---
	BetaAccessCode                                 string `json:"beta_access_code"` // Temporary code for beta access
	FirstName                                      string `json:"first_name"`
	LastName                                       string `json:"last_name"`
	Email                                          string `json:"email"`
	Phone                                          string `json:"phone,omitempty"`
	Country                                        string `json:"country,omitempty"`
	CountryOther                                   string `json:"country_other,omitempty"`
	Timezone                                       string `bson:"timezone" json:"timezone"`
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

func (svc *gatewayFederatedUserRegisterServiceImpl) Execute(
	sessCtx context.Context,
	req *RegisterCustomerRequestIDO,
) error {
	//
	// STEP 1: Sanitization of the input.
	//

	// Defensive Code: For security purposes we need to perform some sanitization on the inputs.
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.ReplaceAll(req.Email, " ", "")
	req.Email = strings.ReplaceAll(req.Email, "\t", "")
	req.Email = strings.TrimSpace(req.Email)

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if req.BetaAccessCode == "" {
		e["beta_access_code"] = "Beta access code is required"
	} else {
		if req.BetaAccessCode != svc.config.App.BetaAccessCode {
			e["beta_access_code"] = "Invalid beta access code"
		}
	}
	if req.FirstName == "" {
		e["first_name"] = "First name is required"
	}
	if req.LastName == "" {
		e["last_name"] = "Last name is required"
	}
	if req.Email == "" {
		e["email"] = "Email is required"
	}
	if len(req.Email) > 255 {
		e["email"] = "Email is too long"
	}
	if req.Phone == "" {
		e["phone"] = "Phone number is required"
	}
	if req.Country == "" {
		e["country"] = "Country is required"
	}
	if req.Timezone == "" {
		e["timezone"] = "Timezone is required"
	}
	if req.AgreeTermsOfService == false {
		e["agree_terms_of_service"] = "Agreeing to terms of service is required and you must agree to the terms before proceeding"
	}
	if req.Module == 0 {
		e["module"] = "Module is required"
	} else {
		// Assuming MonolithModulePaperCloudPropertyEvaluator is the only valid module for now
		if req.Module != int(constants.MonolithModulePaperCloudPropertyEvaluator) {
			e["module"] = "Module is invalid"
		}
	}

	// --- E2EE Related Validation ---
	if req.Salt == "" {
		e["salt"] = "Salt is required"
	}
	if req.PublicKey == "" {
		e["publicKey"] = "Public key is required"
	}
	if req.EncryptedMasterKey == "" {
		e["encryptedMasterKey"] = "Encrypted master key is required"
	}
	if req.EncryptedPrivateKey == "" {
		e["encryptedPrivateKey"] = "Encrypted private key is required"
	}
	if req.EncryptedRecoveryKey == "" {
		e["encryptedRecoveryKey"] = "Encrypted recovery key is required"
	}
	if req.MasterKeyEncryptedWithRecoveryKey == "" {
		e["masterKeyEncryptedWithRecoveryKey"] = "Master key encrypted with recovery key is required"
	}
	if req.VerificationID == "" {
		e["verificationID"] = "Verification ID is required"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the federateduser in our database, else return a `400 Bad Request` error.
	u, err := svc.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		svc.logger.Error("failed getting user by email from database",
			zap.Any("error", err))
		return err
	}
	if u != nil {
		return httperror.NewForBadRequestWithSingleField("email", "Email address already exists")
	}
	// Create our federateduser.
	u, err = svc.createCustomerFederatedUserForRequest(sessCtx, req)
	if err != nil {
		return err
	}

	if err := svc.sendFederatedUserVerificationEmailUseCase.Execute(context.Background(), req.Module, u); err != nil {
		return err
	}

	return nil
}

func (s *gatewayFederatedUserRegisterServiceImpl) createCustomerFederatedUserForRequest(sessCtx context.Context, req *RegisterCustomerRequestIDO) (*dom_user.FederatedUser, error) {

	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	emailVerificationCode, err := random.GenerateSixDigitCode()
	if err != nil {
		return nil, err
	}

	userID := primitive.NewObjectID()
	u := &dom_user.FederatedUser{
		// --- E2EE ---
		Salt:                              req.Salt,
		PublicKey:                         req.PublicKey,
		EncryptedMasterKey:                req.EncryptedMasterKey,
		EncryptedPrivateKey:               req.EncryptedPrivateKey,
		EncryptedRecoveryKey:              req.EncryptedRecoveryKey,
		MasterKeyEncryptedWithRecoveryKey: req.MasterKeyEncryptedWithRecoveryKey,
		VerificationID:                    req.VerificationID,

		// --- The rest of the stuff... ---
		ID:                  userID,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		Name:                fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:         fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:               req.Email,
		Role:                dom_user.FederatedUserRoleIndividual,
		Phone:               req.Phone,
		Country:             req.Country,
		Timezone:            req.Timezone,
		Region:              "",
		City:                "",
		PostalCode:          "",
		AddressLine1:        "",
		AddressLine2:        "",
		AgreeTermsOfService: req.AgreeTermsOfService,
		AgreePromotions:     req.AgreePromotions,
		AgreeToTrackingAcrossThirdPartyAppsAndServices: req.AgreeToTrackingAcrossThirdPartyAppsAndServices,
		CreatedByUserID:         userID,
		CreatedAt:               time.Now(),
		CreatedByName:           fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		CreatedFromIPAddress:    ipAddress,
		ModifiedByUserID:        userID,
		ModifiedAt:              time.Now(),
		ModifiedByName:          fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedFromIPAddress:   ipAddress,
		WasEmailVerified:        false,
		EmailVerificationCode:   fmt.Sprintf("%s", emailVerificationCode),
		EmailVerificationExpiry: time.Now().Add(72 * time.Hour),
		Status:                  dom_user.FederatedUserStatusActive,
		HasShippingAddress:      false,
		ShippingName:            "",
		ShippingPhone:           "",
		ShippingCountry:         "",
		ShippingRegion:          "",
		ShippingCity:            "",
		ShippingPostalCode:      "",
		ShippingAddressLine1:    "",
		ShippingAddressLine2:    "",
	}
	if req.CountryOther != "" {
		u.Country = req.CountryOther
	}
	err = s.userCreateUseCase.Execute(sessCtx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
