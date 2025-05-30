// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/usecase/federateduser/getbyid.go
package federateduser

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type FederatedUserGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_user.FederatedUser, error)
}

type userGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_user.Repository
}

func NewFederatedUserGetByIDUseCase(config *config.Configuration, logger *zap.Logger, repo dom_user.Repository) FederatedUserGetByIDUseCase {
	return &userGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *userGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_user.FederatedUser, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from database.
	//

	return uc.repo.GetByID(ctx, id)
}
