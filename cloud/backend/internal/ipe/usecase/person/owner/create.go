// cloud/backend/internal/ipe/usecase/person/owner/create.go
package owner

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type OwnerCreateUseCase interface {
	Execute(ctx context.Context, owner *dom_person.Owner) (primitive.ObjectID, error)
}

type ownerCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_person.PersonRepository
}

func NewOwnerCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_person.PersonRepository,
) OwnerCreateUseCase {
	return &ownerCreateUseCaseImpl{config, logger, repo}
}

func (uc *ownerCreateUseCaseImpl) Execute(ctx context.Context, owner *dom_person.Owner) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if owner == nil {
		e["owner"] = "Owner is required"
	} else {
		if owner.PersonName == "" {
			e["person_name"] = "Person name is required"
		}
		// Additional validations
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating owner creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.SaveOwner(ctx, owner)
}
