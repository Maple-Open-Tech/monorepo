package gateway

import (
	"context"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/database/mongodbcache"
)

type GatewayLogoutService interface {
	Execute(ctx context.Context) error
}

type gatewayLogoutServiceImpl struct {
	cache mongodbcache.Cacher
}

func NewGatewayLogoutService(
	cach mongodbcache.Cacher,
) GatewayLogoutService {
	return &gatewayLogoutServiceImpl{cach}
}

func (s *gatewayLogoutServiceImpl) Execute(ctx context.Context) error {
	// Extract from our session the following data.
	sessionID, ok := ctx.Value(constants.SessionID).(string)
	if !ok {
		return httperror.NewForBadRequestWithSingleField("session_id", "not logged in")
	}

	if err := s.cache.Delete(ctx, sessionID); err != nil {
		return err
	}
	return nil
}
