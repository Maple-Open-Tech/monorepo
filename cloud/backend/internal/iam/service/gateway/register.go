package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
	uc_emailer "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/emailer"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/random"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	sstring "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/securestring"
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
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.FederatedUserGetByEmailUseCase,
	uc2 uc_user.FederatedUserCreateUseCase,
	uc3 uc_user.FederatedUserUpdateUseCase,
	uc4 uc_emailer.SendFederatedUserVerificationEmailUseCase,
) GatewayFederatedUserRegisterService {
	return &gatewayFederatedUserRegisterServiceImpl{cfg, pp, cach, jwtp, uc1, uc2, uc3, uc4}
}

type RegisterCustomerRequestIDO struct {
	BetaAccessCode                                 string `json:"beta_access_code"` // Temporary code for beta access
	FirstName                                      string `json:"first_name"`
	LastName                                       string `json:"last_name"`
	Email                                          string `json:"email"`
	Password                                       string `json:"password"`
	PasswordConfirm                                string `json:"password_confirm"`
	Phone                                          string `json:"phone,omitempty"`
	Country                                        string `json:"country,omitempty"`
	CountryOther                                   string `json:"country_other,omitempty"`
	Timezone                                       string `bson:"timezone" json:"timezone"`
	AgreeTermsOfService                            bool   `json:"agree_terms_of_service,omitempty"`
	AgreePromotions                                bool   `json:"agree_promotions,omitempty"`
	AgreeToTrackingAcrossThirdPartyAppsAndServices bool   `json:"agree_to_tracking_across_third_party_apps_and_services,omitempty"`

	// Module refers to which module the user is registering for.
	Module int `json:"module,omitempty"`
}

type RegisterCustomerResponseIDO struct {
	FederatedUser          *domain.FederatedUser `json:"federateduser"`
	AccessToken            string                `json:"access_token"`
	AccessTokenExpiryTime  time.Time             `json:"access_token_expiry_time"`
	RefreshToken           string                `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time             `json:"refresh_token_expiry_time"`
}

func (s *gatewayFederatedUserRegisterServiceImpl) Execute(
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
	req.Password = strings.ReplaceAll(req.Password, " ", "")
	req.Password = strings.ReplaceAll(req.Password, "\t", "")
	req.Password = strings.TrimSpace(req.Password)
	// password, err := sstring.NewSecureString(unsecurePassword)

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if req.BetaAccessCode == "" {
		e["beta_access_code"] = "Beta access code is required"
	} else {
		if req.BetaAccessCode != s.config.App.BetaAccessCode {
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
	if req.Password == "" {
		e["password"] = "Password is required"
	}
	if req.PasswordConfirm == "" {
		e["password_confirm"] = "Password confirm is required"
	}
	if req.PasswordConfirm != req.Password {
		e["password"] = "Password does not match"
		e["password_confirm"] = "Password does not match"
	}
	if req.Phone == "" {
		e["phone"] = "Phone confirm is required"
	}
	if req.Country == "" {
		e["country"] = "Country is required"
	} else {
		if req.Country == "Other" && req.CountryOther == "" {
			e["country_other"] = "Specify country is required"
		}
	}
	if req.Timezone == "" {
		e["timezone"] = "Password confirm is required"
	}
	if req.AgreeTermsOfService == false {
		e["agree_terms_of_service"] = "Agreeing to terms of service is required and you must agree to the terms before proceeding"
	}
	if req.Module == 0 {
		e["module"] = "Module is required"
	} else {
		if req.Module != int(constants.MonolithModuleIncomePropertyEvaluator) {
			e["module"] = "Module is invalid"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the federateduser in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		return err
	}
	if u != nil {
		return httperror.NewForBadRequestWithSingleField("email", "Email address already exists")
	}

	// Create our federateduser.
	u, err = s.createCustomerFederatedUserForRequest(sessCtx, req)
	if err != nil {
		return err
	}

	if err := s.sendFederatedUserVerificationEmailUseCase.Execute(context.Background(), req.Module, u); err != nil {
		return err
	}

	return nil
}

func (s *gatewayFederatedUserRegisterServiceImpl) createCustomerFederatedUserForRequest(sessCtx context.Context, req *RegisterCustomerRequestIDO) (*domain.FederatedUser, error) {

	password, err := sstring.NewSecureString(req.Password)
	if err != nil {
		return nil, err
	}
	defer password.Wipe()

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(password)
	if err != nil {
		return nil, err
	}

	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	emailVerificationCode, err := random.GenerateSixDigitCode()
	if err != nil {
		return nil, err
	}

	userID := primitive.NewObjectID()
	u := &domain.FederatedUser{
		ID:                    userID,
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Name:                  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:           fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                 req.Email,
		PasswordHash:          passwordHash,
		PasswordHashAlgorithm: s.passwordProvider.AlgorithmName(),
		Role:                  domain.FederatedUserRoleIndividual,
		Phone:                 req.Phone,
		Country:               req.Country,
		Timezone:              req.Timezone,
		Region:                "",
		City:                  "",
		PostalCode:            "",
		AddressLine1:          "",
		AddressLine2:          "",
		AgreeTermsOfService:   req.AgreeTermsOfService,
		AgreePromotions:       req.AgreePromotions,
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
		Status:                  domain.FederatedUserStatusActive,
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
