// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/user/getbyid.go
package user

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/user"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type UserGetByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*dom_user.User, error)
}

type userGetByIDUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_user.Repository
}

func NewUserGetByIDUseCase(config *config.Configuration, logger *zap.Logger, repo dom_user.Repository) UserGetByIDUseCase {
	return &userGetByIDUseCaseImpl{config, logger, repo}
}

func (uc *userGetByIDUseCaseImpl) Execute(ctx context.Context, id primitive.ObjectID) (*dom_user.User, error) {
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
