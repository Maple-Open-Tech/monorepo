package gateway

import (
	"context"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayLogoutService interface {
	Execute(ctx context.Context) error
}

type gatewayLogoutServiceImpl struct {
	logger *zap.Logger
	cache  mongodbcache.Cacher
}

func NewGatewayLogoutService(
	logger *zap.Logger,
	cach mongodbcache.Cacher,
) GatewayLogoutService {
	return &gatewayLogoutServiceImpl{logger, cach}
}

func (s *gatewayLogoutServiceImpl) Execute(ctx context.Context) error {
	// Extract from our session the following data.
	sessionID, ok := ctx.Value(constants.SessionID).(string)
	if !ok {
		s.logger.Warn("loggout could not happen - no session in mongo-cache")
		return httperror.NewForBadRequestWithSingleField("session_id", "not logged in")
	}

	if err := s.cache.Delete(ctx, sessionID); err != nil {
		s.logger.Error("cache delete error", zap.Any("err", err))
		return err
	}
	return nil
}
