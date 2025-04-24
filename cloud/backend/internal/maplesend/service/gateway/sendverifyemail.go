package gateway

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/mongo"

	uc_emailer "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/emailer"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/usecase/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type GatewaySendVerifyEmailService interface {
	Execute(sessCtx mongo.SessionContext, req *GatewaySendVerifyEmailRequestIDO) error
}

type gatewaySendVerifyEmailServiceImpl struct {
	logger                           *zap.Logger
	userGetByEmailUseCase            uc_user.UserGetByEmailUseCase
	sendUserVerificationEmailUseCase uc_emailer.SendUserVerificationEmailUseCase
}

func NewGatewaySendVerifyEmailService(
	logger *zap.Logger,
	uc1 uc_user.UserGetByEmailUseCase,
	uc2 uc_emailer.SendUserVerificationEmailUseCase,
) GatewaySendVerifyEmailService {
	return &gatewaySendVerifyEmailServiceImpl{logger, uc1, uc2}
}

type GatewaySendVerifyEmailRequestIDO struct {
	Email string `json:"email"`
}

func (s *gatewaySendVerifyEmailServiceImpl) Execute(sessCtx mongo.SessionContext, req *GatewaySendVerifyEmailRequestIDO) error {
	// Extract from our session the following data.
	// sessionID := sessCtx.Value(constants.SessionID).(string)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		s.logger.Error("database error", zap.Any("err", err))
		return err
	}
	if u == nil {
		s.logger.Warn("user does not exist for email error")
		return httperror.NewForBadRequestWithSingleField("email", "does not exist")
	}

	if err := s.sendUserVerificationEmailUseCase.Execute(context.Background(), u); err != nil {
		s.logger.Error("failed sending verification email with error", zap.Any("err", err))
		// Skip any error handling...
	}

	return nil
}
