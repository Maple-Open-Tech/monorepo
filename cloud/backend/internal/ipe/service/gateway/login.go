package gateway

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/user"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/user"
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
	logger                *zap.Logger
	passwordProvider      password.Provider
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	userGetByEmailUseCase uc_user.UserGetByEmailUseCase
	userUpdateUseCase     uc_user.UserUpdateUseCase
}

func NewGatewayLoginService(
	logger *zap.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.UserGetByEmailUseCase,
	uc2 uc_user.UserUpdateUseCase,
) GatewayLoginService {
	return &gatewayLoginServiceImpl{logger, pp, cach, jwtp, uc1, uc2}
}

type GatewayLoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GatewayLoginResponseIDO struct {
	User                   *domain.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
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
		s.logger.Warn("Failed validation login",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		s.logger.Error("database error",
			zap.String("email", req.Email),
			zap.Any("err", err))
		return nil, err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error",
			zap.String("email", req.Email))
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address does not exist")
	}

	s.logger.Debug("attempting to confirm correct password submission for the existing user...",
		zap.String("email", req.Email))

	securePassword, err := sstring.NewSecureString(req.Password)
	if err != nil {
		s.logger.Error("database error",
			zap.String("email", req.Email),
			zap.Any("err", err))
		return nil, err
	}
	defer securePassword.Wipe()

	s.logger.Debug("attempting to compare password hashes...",
		zap.String("email", req.Email))

	// Verify the inputted password and hashed password match.
	passwordMatch, _ := s.passwordProvider.ComparePasswordAndHash(securePassword, u.PasswordHash)
	if passwordMatch == false {
		s.logger.Warn("password check validation error",
			zap.String("email", req.Email))
		return nil, httperror.NewForBadRequestWithSingleField("password", "Password does not match with record")
	}

	s.logger.Debug("attempting to confirm existing user has a verified email address...",
		zap.String("email", req.Email))

	// Enforce the verification code of the email.
	if u.WasEmailVerified == false {
		s.logger.Warn("email verification validation error",
			zap.String("email", req.Email))
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address was not verified")
	}

	s.logger.Debug("login confirmed correct password and verified email for existing user...",
		zap.String("email", req.Email))

	// // Enforce 2FA if enabled.
	if u.OTPEnabled {
		// We need to reset the `otp_validated` status to be false to force
		// the user to use their `totp authenticator` application.
		u.OTPValidated = false
		u.ModifiedAt = time.Now()
		if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
			s.logger.Error("failed updating user during login",
				zap.String("email", req.Email),
				zap.Any("err", err))
			return nil, err
		}
	}

	return s.loginWithUser(sessCtx, u)
}

func (s *gatewayLoginServiceImpl) loginWithUser(sessCtx context.Context, u *domain.User) (*GatewayLoginResponseIDO, error) {
	uBin, err := json.Marshal(u)
	if err != nil {
		s.logger.Error("marshalling error", zap.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 5 * time.Minute     // 5 minutes
	rtExpiry := 14 * 24 * time.Hour // 1 week

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

	// For debugging purposes we want to print the wallet address.
	s.logger.Debug("login successfull")

	// Return our auth keys.
	return &GatewayLoginResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}
