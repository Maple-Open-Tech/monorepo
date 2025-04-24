package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	uc_emailer "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/emailer"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/random"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/password"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayForgotPasswordService interface {
	Execute(sessCtx context.Context, req *GatewayForgotPasswordRequestIDO) (*GatewayForgotPasswordResponseIDO, error)
}

type gatewayForgotPasswordServiceImpl struct {
	logger                            *zap.Logger
	passwordProvider                  password.Provider
	cache                             mongodbcache.Cacher
	jwtProvider                       jwt.Provider
	userGetByEmailUseCase             uc_user.UserGetByEmailUseCase
	userUpdateUseCase                 uc_user.UserUpdateUseCase
	sendUserPasswordResetEmailUseCase uc_emailer.SendUserPasswordResetEmailUseCase
}

func NewGatewayForgotPasswordService(
	logger *zap.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 uc_user.UserGetByEmailUseCase,
	uc2 uc_user.UserUpdateUseCase,
	uc3 uc_emailer.SendUserPasswordResetEmailUseCase,
) GatewayForgotPasswordService {
	return &gatewayForgotPasswordServiceImpl{logger, pp, cach, jwtp, uc1, uc2, uc3}
}

type GatewayForgotPasswordRequestIDO struct {
	Email string `json:"email"`
}

type GatewayForgotPasswordResponseIDO struct {
	Message string `json:"message"`
}

func (s *gatewayForgotPasswordServiceImpl) Execute(sessCtx context.Context, req *GatewayForgotPasswordRequestIDO) (*GatewayForgotPasswordResponseIDO, error) {
	//
	// STEP 1: Sanization of input.
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
	if req.Email == "" {
		e["email"] = "Email address is required"
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
		s.logger.Error("database error", zap.Any("err", err))
		return nil, err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address does not exist")
	}

	//
	// STEP 4:
	//

	passwordResetVerificationCode, err := random.GenerateSixDigitCode()
	if err != nil {
		s.logger.Error("generating email verification code error", zap.Any("error", err))
		return nil, err
	}

	u.PasswordResetVerificationCode = fmt.Sprintf("%s", passwordResetVerificationCode)
	u.PasswordResetVerificationExpiry = time.Now().Add(5 * time.Minute)
	u.ModifiedAt = time.Now()
	u.ModifiedByName = u.Name
	err = s.userUpdateUseCase.Execute(sessCtx, u)
	if err != nil {
		s.logger.Error("database update error", zap.Any("error", err))
		return nil, err
	}

	//
	// STEP 5: Send email
	//

	if err := s.sendUserPasswordResetEmailUseCase.Execute(sessCtx, u); err != nil {
		s.logger.Error("failed sending verification email with error", zap.Any("err", err))
		// Skip any error handling...
	}

	//
	// STEP X: Done
	//

	// Return our auth keys.
	return &GatewayForgotPasswordResponseIDO{
		Message: "Password reset email has been sent",
	}, nil
}
