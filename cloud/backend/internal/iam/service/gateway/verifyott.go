// cloud/backend/internal/iam/service/gateway/verifyott.go
package gateway

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/nacl/box"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

// Data structures for OTT verification
type GatewayVerifyLoginOTTRequestIDO struct {
	Email string `json:"email"`
	OTT   string `json:"ott"`
}

type GatewayVerifyLoginOTTResponseIDO struct {
	Salt                string `json:"salt"`
	PublicKey           string `json:"publicKey"`
	EncryptedMasterKey  string `json:"encryptedMasterKey"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey"`
	EncryptedChallenge  string `json:"encryptedChallenge"`
	ChallengeID         string `json:"challengeId"`
}

// ChallengeData structure to be stored in cache
type ChallengeData struct {
	Email           string    `json:"email"`
	ChallengeID     string    `json:"challenge_id"`
	Challenge       string    `json:"challenge"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	IsVerified      bool      `json:"is_verified"`
	FederatedUserID string    `json:"federated_user_id"`
}

// Service interface for OTT verification
type GatewayVerifyLoginOTTService interface {
	Execute(sessCtx context.Context, req *GatewayVerifyLoginOTTRequestIDO) (*GatewayVerifyLoginOTTResponseIDO, error)
}

// Implementation of OTT verification service
type gatewayVerifyLoginOTTServiceImpl struct {
	config                *config.Configuration
	logger                *zap.Logger
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	userGetByEmailUseCase uc_user.FederatedUserGetByEmailUseCase
}

func NewGatewayVerifyLoginOTTService(
	config *config.Configuration,
	logger *zap.Logger,
	cache mongodbcache.Cacher,
	jwtProvider jwt.Provider,
	userGetByEmailUseCase uc_user.FederatedUserGetByEmailUseCase,
) GatewayVerifyLoginOTTService {
	return &gatewayVerifyLoginOTTServiceImpl{
		config:                config,
		logger:                logger,
		cache:                 cache,
		jwtProvider:           jwtProvider,
		userGetByEmailUseCase: userGetByEmailUseCase,
	}
}

func (s *gatewayVerifyLoginOTTServiceImpl) Execute(sessCtx context.Context, req *GatewayVerifyLoginOTTRequestIDO) (*GatewayVerifyLoginOTTResponseIDO, error) {
	// Validate input
	e := make(map[string]string)
	if req.Email == "" {
		e["email"] = "Email address is required"
	}
	if req.OTT == "" {
		e["ott"] = "Verification code is required"
	}
	if len(e) != 0 {
		return nil, httperror.NewForBadRequest(&e)
	}

	// Sanitize input
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.ReplaceAll(req.Email, " ", "")
	req.OTT = strings.TrimSpace(req.OTT)

	// Retrieve OTT data from cache
	cacheKey := fmt.Sprintf("login_ott:%s", req.Email)
	ottDataJSON, err := s.cache.Get(sessCtx, cacheKey)
	if err != nil {
		s.logger.Error("Failed to retrieve OTT data", zap.Error(err))
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Invalid or expired verification code")
	}

	if ottDataJSON == nil {
		s.logger.Error("OTT data not found in cache")
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Invalid or expired verification code")
	}

	// Unmarshal the data from JSON
	var ottData LoginOTTData
	if err := json.Unmarshal(ottDataJSON, &ottData); err != nil {
		s.logger.Error("Failed to unmarshal OTT data", zap.Error(err))
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Invalid verification code")
	}

	// Verify OTT
	if ottData.OTT != req.OTT {
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Invalid verification code")
	}

	// Check expiry
	if time.Now().After(ottData.ExpiresAt) {
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Verification code has expired")
	}

	// Check if already verified
	if ottData.IsVerified {
		return nil, httperror.NewForBadRequestWithSingleField("ott", "Verification code has already been used")
	}

	// Get user from database
	user, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address does not exist")
	}

	// Generate a challenge for final verification
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		s.logger.Error("Failed to generate challenge", zap.Error(err))
		return nil, fmt.Errorf("failed to process login: %w", err)
	}

	// Base64 encode the challenge for storage
	challengeBase64 := base64.StdEncoding.EncodeToString(challenge)

	// Generate a unique challenge ID
	challengeID := uuid.New().String()

	// Store challenge in cache
	challengeData := ChallengeData{
		Email:           req.Email,
		ChallengeID:     challengeID,
		Challenge:       challengeBase64,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(5 * time.Minute), // Challenge valid for 5 minutes
		IsVerified:      false,
		FederatedUserID: user.ID.Hex(),
	}

	// Generate a unique cache key for this challenge
	challengeCacheKey := fmt.Sprintf("login_challenge:%s", challengeID)

	// Marshal the challenge data to JSON
	challengeDataJSON, err := json.Marshal(challengeData)
	if err != nil {
		s.logger.Error("Failed to marshal challenge data", zap.Error(err))
		return nil, fmt.Errorf("failed to process login verification: %w", err)
	}

	// Store in cache with expiry
	if err := s.cache.SetWithExpiry(sessCtx, challengeCacheKey, challengeDataJSON, 5*time.Minute); err != nil {
		s.logger.Error("Failed to store challenge in cache", zap.Error(err))
		return nil, fmt.Errorf("failed to process login verification: %w", err)
	}

	// Mark OTT as verified
	ottData.IsVerified = true
	ottData.ChallengeID = challengeID

	// Marshal the updated OTT data to JSON
	updatedOTTDataJSON, err := json.Marshal(ottData)
	if err != nil {
		s.logger.Error("Failed to marshal updated OTT data", zap.Error(err))
		// Continue anyway, as the challenge is already stored
	} else {
		if err := s.cache.SetWithExpiry(sessCtx, cacheKey, updatedOTTDataJSON, 10*time.Minute); err != nil {
			s.logger.Error("Failed to update OTT in cache", zap.Error(err))
			// Continue anyway, as the challenge is already stored
		}
	}

	encryptedChallenge, err := getEncryptedChallenge(challenge, user)
	if err != nil {
		s.logger.Error("Failed to encrypt challenge", zap.Error(err))
		return nil, fmt.Errorf("failed to process login: %w", err)
	}

	// Return encrypted keys and challenge for client-side password verification
	return &GatewayVerifyLoginOTTResponseIDO{
		Salt:                user.Salt,
		PublicKey:           user.PublicKey,
		EncryptedMasterKey:  user.EncryptedMasterKey,
		EncryptedPrivateKey: user.EncryptedPrivateKey,
		EncryptedChallenge:  encryptedChallenge,
		ChallengeID:         challengeID,
	}, nil
}

// getEncryptedChallenge encrypts the challenge with the user's master key
func getEncryptedChallenge(challenge []byte, user *domain.FederatedUser) (string, error) {
	// Decode the user's public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(user.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode user public key: %w", err)
	}

	// Ensure we have the right length for NaCl box
	if len(publicKeyBytes) != 32 {
		return "", fmt.Errorf("invalid public key length: got %d, want 32", len(publicKeyBytes))
	}

	// Generate a random nonce
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create a ephemeral keypair for this encryption
	ephemeralPub, ephemeralPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate ephemeral keypair: %w", err)
	}

	// Convert the user's public key to the expected format
	var userPubKey [32]byte
	copy(userPubKey[:], publicKeyBytes)

	// Encrypt the challenge with box.Seal using the user's public key
	encrypted := box.Seal(nonce[:], challenge, &nonce, &userPubKey, ephemeralPriv)

	// Prepend the ephemeral public key to the encrypted data
	result := make([]byte, len(encrypted)+32)
	copy(result[:32], ephemeralPub[:])
	copy(result[32:], encrypted)

	// Return base64 encoded result
	return base64.StdEncoding.EncodeToString(result), nil
}
