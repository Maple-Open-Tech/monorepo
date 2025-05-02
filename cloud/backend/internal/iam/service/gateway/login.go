package gateway

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	sstring "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/securestring"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayLoginService interface {
	Execute(sessCtx context.Context, req *GatewayLoginRequestIDO) (*GatewayLoginResponseIDO, error)
}

type gatewayLoginServiceImpl struct {
	passwordProvider      password.Provider
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	userGetByEmailUseCase uc_user.FederatedUserGetByEmailUseCase
	userUpdateUseCase     uc_user.FederatedUserUpdateUseCase
}

func NewGatewayLoginService(
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.FederatedUserGetByEmailUseCase,
	uc2 uc_user.FederatedUserUpdateUseCase,
) GatewayLoginService {
	return &gatewayLoginServiceImpl{pp, cach, jwtp, uc1, uc2}
}

type GatewayLoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GatewayLoginResponseIDO struct {
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}

func (s *gatewayLoginServiceImpl) Execute(sessCtx context.Context, req *GatewayLoginRequestIDO) (*GatewayLoginResponseIDO, error) {
	//
	// STEP 1: Sanization of input.
	//

	// Defensive Code: For security purposes we need to perform some sanitization on the inputs.
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.ReplaceAll(req.Email, " ", "")
	req.Email = strings.ReplaceAll(req.Email, "\t", "")
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")
	req.Password = strings.ReplaceAll(req.Password, "\t", "")
	req.Password = strings.TrimSpace(req.Password)

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if req.Email == "" {
		e["email"] = "Email address is required"
	}
	if req.Password == "" {
		e["password"] = "Password is required"
	}

	if len(e) != 0 {
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the federateduser in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address does not exist")
	}

	securePassword, err := sstring.NewSecureString(req.Password)
	if err != nil {
		return nil, err
	}
	defer securePassword.Wipe()

	// Verify the inputted password and hashed password match.
	passwordMatch, _ := s.passwordProvider.ComparePasswordAndHash(securePassword, u.PasswordHash)
	if passwordMatch == false {
		return nil, httperror.NewForBadRequestWithSingleField("password", "Password does not match with record")
	}

	// Enforce the verification code of the email.
	if u.WasEmailVerified == false {
		return nil, httperror.NewForBadRequestWithSingleField("email", "Your email address has not been verified. Please check your inbox for the verification email or use the 'Resend Verification Email' option.")
	}

	// // Enforce 2FA if enabled.
	if u.OTPEnabled {
		// We need to reset the `otp_validated` status to be false to force
		// the federateduser to use their `totp authenticator` application.
		u.OTPValidated = false
		u.ModifiedAt = time.Now()
		if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
			return nil, err
		}
	}

	return s.loginWithFederatedUser(sessCtx, u)
}

func (s *gatewayLoginServiceImpl) loginWithFederatedUser(sessCtx context.Context, u *domain.FederatedUser) (*GatewayLoginResponseIDO, error) {
	uBin, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 5 * time.Minute     // 5 minutes
	rtExpiry := 14 * 24 * time.Hour // 1 week

	// Start our session using an access and refresh token.
	sessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, sessionUUID, uBin, rtExpiry)
	if err != nil {
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		return nil, err
	}

	// Return our auth keys.
	return &GatewayLoginResponseIDO{
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}
