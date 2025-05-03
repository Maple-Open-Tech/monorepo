// cloud/backend/internal/ipe/usecase/person/presenter/create.go
package presenter

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type PresenterCreateUseCase interface {
	Execute(ctx context.Context, presenter *dom_person.Presenter) (primitive.ObjectID, error)
}

type presenterCreateUseCaseImpl struct {
	config *config.Configuration
	logger *zap.Logger
	repo   dom_person.PersonRepository
}

func NewPresenterCreateUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repo dom_person.PersonRepository,
) PresenterCreateUseCase {
	return &presenterCreateUseCaseImpl{config, logger, repo}
}

func (uc *presenterCreateUseCaseImpl) Execute(ctx context.Context, presenter *dom_person.Presenter) (primitive.ObjectID, error) {
	//
	// STEP 1: Validation.
	//
	e := make(map[string]string)
	if presenter == nil {
		e["presenter"] = "Presenter is required"
	} else {
		if presenter.PersonName == "" {
			e["person_name"] = "Person name is required"
		}
		// Additional validations
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating presenter creation",
			zap.Any("error", e))
		return primitive.NilObjectID, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//
	return uc.repo.SavePresenter(ctx, presenter)
}
