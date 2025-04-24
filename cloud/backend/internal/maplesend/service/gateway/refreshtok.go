package gateway

import (
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/domain/user"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayRefreshTokenService interface {
	Execute(
		sessCtx mongo.SessionContext,
		req *GatewayRefreshTokenRequestIDO,
	) (*GatewayRefreshTokenResponseIDO, error)
}

type gatewayRefreshTokenServiceImpl struct {
	logger                *zap.Logger
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	userGetByEmailUseCase uc_user.UserGetByEmailUseCase
}

func NewGatewayRefreshTokenService(
	logger *zap.Logger,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.UserGetByEmailUseCase,
) GatewayRefreshTokenService {
	return &gatewayRefreshTokenServiceImpl{logger, cach, jwtp, uc1}
}

type GatewayRefreshTokenRequestIDO struct {
	Value string `json:"value"`
}

// GatewayRefreshTokenResponseIDO struct used to represent the system's response when the `login` POST request was a success.
type GatewayRefreshTokenResponseIDO struct {
	Email                  string    `json:"username"`
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryDate  time.Time `json:"access_token_expiry_date"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryDate time.Time `json:"refresh_token_expiry_date"`
}

func (s *gatewayRefreshTokenServiceImpl) Execute(
	sessCtx mongo.SessionContext,
	req *GatewayRefreshTokenRequestIDO,
) (*GatewayRefreshTokenResponseIDO, error) {
	////
	//// Extract the `sessionID` so we can process it.
	////

	sessionID, err := s.jwtProvider.ProcessJWTToken(req.Value)
	if err != nil {
		s.logger.Warn("process jwt refresh token does not exist", zap.String("value", req.Value))
		err := errors.New("jwt refresh token failed")
		return nil, err
	}

	////
	//// Lookup in our in-memory the user record for the `sessionID` or error.
	////

	uBin, err := s.cache.Get(sessCtx, sessionID)
	if err != nil {
		s.logger.Error("in-memory set error", zap.Any("err", err))
		return nil, err
	}

	var u *domain.User
	err = json.Unmarshal(uBin, &u)
	if err != nil {
		s.logger.Error("unmarshal error", zap.Any("err", err))
		return nil, err
	}

	////
	//// Generate new access and refresh tokens and return them.
	////

	// Set expiry duration.
	atExpiry := 5 * time.Minute     // 5 minutes
	rtExpiry := 14 * 24 * time.Hour // 1 week

	// Start our session using an access and refresh token.
	newSessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, newSessionUUID, uBin, rtExpiry)
	if err != nil {
		s.logger.Error("cache set with expiry error", zap.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(newSessionUUID, atExpiry, rtExpiry)
	if err != nil {
		s.logger.Error("jwt generate pairs error", zap.Any("err", err))
		return nil, err
	}

	ido := &GatewayRefreshTokenResponseIDO{
		Email:                  u.Email,
		AccessToken:            accessToken,
		AccessTokenExpiryDate:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryDate: refreshTokenExpiry,
	}

	// Return our auth keys.
	return ido, nil
}
