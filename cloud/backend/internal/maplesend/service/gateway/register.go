package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/domain/user"
	uc_emailer "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/emailer"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/random"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	sstring "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/securestring"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayUserRegisterService interface {
	Execute(
		sessCtx mongo.SessionContext,
		req *RegisterCustomerRequestIDO,
	) (*RegisterCustomerResponseIDO, error)
}

type gatewayUserRegisterServiceImpl struct {
	config                           *config.Configuration
	logger                           *zap.Logger
	passwordProvider                 password.Provider
	cache                            mongodbcache.Cacher
	jwtProvider                      jwt.Provider
	userGetByEmailUseCase            uc_user.UserGetByEmailUseCase
	userCreateUseCase                uc_user.UserCreateUseCase
	userUpdateUseCase                uc_user.UserUpdateUseCase
	sendUserVerificationEmailUseCase uc_emailer.SendUserVerificationEmailUseCase
}

func NewGatewayUserRegisterService(
	cfg *config.Configuration,
	logger *zap.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.UserGetByEmailUseCase,
	uc2 uc_user.UserCreateUseCase,
	uc3 uc_user.UserUpdateUseCase,
	uc4 uc_emailer.SendUserVerificationEmailUseCase,
) GatewayUserRegisterService {
	return &gatewayUserRegisterServiceImpl{cfg, logger, pp, cach, jwtp, uc1, uc2, uc3, uc4}
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
}

type RegisterCustomerResponseIDO struct {
	User                   *domain.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

func (s *gatewayUserRegisterServiceImpl) Execute(
	sessCtx mongo.SessionContext,
	req *RegisterCustomerRequestIDO,
) (*RegisterCustomerResponseIDO, error) {
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
	// if err != nil {
	// 	s.logger.Error("secure string error", zap.Any("err", err))
	// 	return nil, err
	// }

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
	if len(e) != 0 {
		s.logger.Warn("Failed validation register",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		s.logger.Error("database error", zap.Any("err", err))
		return nil, err
	}
	if u != nil {
		s.logger.Warn("user already exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address already exists")
	}

	// Create our user.
	u, err = s.createCustomerUserForRequest(sessCtx, req)
	if err != nil {
		s.logger.Error("failed creating customer user error", zap.Any("err", err))
		return nil, err
	}

	// Send our verification email.
	if err := s.sendUserVerificationEmailUseCase.Execute(context.Background(), u); err != nil {
		s.logger.Error("failed sending verification email with error", zap.Any("err", err))
		// Skip any error handling...
	}

	return s.registerWithUser(sessCtx, u)
}

func (s *gatewayUserRegisterServiceImpl) registerWithUser(sessCtx mongo.SessionContext, u *domain.User) (*RegisterCustomerResponseIDO, error) {
	uBin, err := json.Marshal(u)
	if err != nil {
		s.logger.Error("marshalling error", zap.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	sessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, sessionUUID, uBin, rtExpiry)
	if err != nil {
		s.logger.Error("cache set with expiry error", zap.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		s.logger.Error("jwt generate pairs error", zap.Any("err", err))
		return nil, err
	}

	// Return our auth keys.
	return &RegisterCustomerResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}

func (s *gatewayUserRegisterServiceImpl) createCustomerUserForRequest(sessCtx mongo.SessionContext, req *RegisterCustomerRequestIDO) (*domain.User, error) {

	password, err := sstring.NewSecureString(req.Password)
	if err != nil {
		s.logger.Error("password securing error", zap.Any("err", err))
		return nil, err
	}
	defer password.Wipe()

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(password)
	if err != nil {
		s.logger.Error("hashing error", zap.Any("error", err))
		return nil, err
	}

	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	emailVerificationCode, err := random.GenerateSixDigitCode()
	if err != nil {
		s.logger.Error("generating email verification code error", zap.Any("error", err))
		return nil, err
	}

	userID := primitive.NewObjectID()
	u := &domain.User{
		ID:                    userID,
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Name:                  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:           fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                 req.Email,
		PasswordHash:          passwordHash,
		PasswordHashAlgorithm: s.passwordProvider.AlgorithmName(),
		Role:                  domain.UserRoleIndividual,
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
		CreatedByUserID:            userID,
		CreatedAt:                  time.Now(),
		CreatedByName:              fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		CreatedFromIPAddress:       ipAddress,
		ModifiedByUserID:           userID,
		ModifiedAt:                 time.Now(),
		ModifiedByName:             fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedFromIPAddress:      ipAddress,
		WasEmailVerified:           false,
		EmailVerificationCode:      fmt.Sprintf("%s", emailVerificationCode),
		EmailVerificationExpiry:    time.Now().Add(72 * time.Hour),
		Status:                     domain.UserStatusActive,
		HasShippingAddress:         false,
		ShippingName:               "",
		ShippingPhone:              "",
		ShippingCountry:            "",
		ShippingRegion:             "",
		ShippingCity:               "",
		ShippingPostalCode:         "",
		ShippingAddressLine1:       "",
		ShippingAddressLine2:       "",
		PaymentProcessorName:       "",
		PaymentProcessorCustomerID: "",
		// HowDidYouHearAboutUs:                  "",
		// HowDidYouHearAboutUsOther:             "",
		// HowLongCollectingComicBooksForGrading: "",
		// HasPreviouslySubmittedComicBookForGrading:       "",
		// HasOwnedGradedComicBooks:                        "",
		// HasRegularComicBookShop:                         "",
		// HasPreviouslyPurchasedFromAuctionSite:           "",
		// HasPreviouslyPurchasedFromFacebookMarketplace:   "",
		// HasRegularlyAttendedComicConsOrCollectibleShows: "",
		ProfileVerificationStatus: domain.UserProfileVerificationStatusUnverified,
	}
	if req.CountryOther != "" {
		u.Country = req.CountryOther
	}
	err = s.userCreateUseCase.Execute(sessCtx, u)
	if err != nil {
		s.logger.Error("database create error", zap.Any("error", err))
		return nil, err
	}
	s.logger.Info("Customer user created.",
		zap.Any("_id", u.ID),
		zap.String("full_name", u.Name),
		zap.String("email", u.Email),
		zap.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		zap.String("password_hash", u.PasswordHash))

	return u, nil
}
