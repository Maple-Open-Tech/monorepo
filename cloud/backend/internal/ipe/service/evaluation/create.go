// cloud/backend/internal/ipe/service/evaluation/create.go
package evaluation

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	uc_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateEvaluationRequestDTO struct {
	PropertyID                                    primitive.ObjectID `json:"propertyId"`
	ClientID                                      primitive.ObjectID `json:"clientId"`
	PresenterID                                   primitive.ObjectID `json:"presenterId"`
	OwnerID                                       primitive.ObjectID `json:"ownerId"`
	ShouldDisplayTitleSection                     bool               `json:"shouldDisplayTitleSection"`
	ShouldDisplayExecutiveSummarySection          bool               `json:"shouldDisplayExecutiveSummarySection"`
	ShouldDisplayLocationSection                  bool               `json:"shouldDisplayLocationSection"`
	ShouldDisplayNeighbourhoodSection             bool               `json:"shouldDisplayNeighbourhoodSection"`
	ShouldDisplayExteriorSection                  bool               `json:"shouldDisplayExteriorSection"`
	ShouldDisplayBuildingSection                  bool               `json:"shouldDisplayBuildingSection"`
	ShouldDisplayCommercialInteriorSection        bool               `json:"shouldDisplayCommercialInteriorSection"`
	ShouldDisplayCommercialInteriorDetailsSection bool               `json:"shouldDisplayCommercialInteriorDetailsSection"`
	ShouldDisplayResidentialInteriorSection       bool               `json:"shouldDisplayResidentialInteriorSection"`
	ShouldDisplayLegalSection                     bool               `json:"shouldDisplayLegalSection"`
	ShouldDisplayFinancialSection                 bool               `json:"shouldDisplayFinancialSection"`
	NeighbourhoodID                               primitive.ObjectID `json:"neighbourhoodId"`
}

type CreateEvaluationResponseDTO struct {
	ID                                            primitive.ObjectID `json:"id"`
	PropertyID                                    primitive.ObjectID `json:"propertyId"`
	ClientID                                      primitive.ObjectID `json:"clientId"`
	PresenterID                                   primitive.ObjectID `json:"presenterId"`
	OwnerID                                       primitive.ObjectID `json:"ownerId"`
	ShouldDisplayTitleSection                     bool               `json:"shouldDisplayTitleSection"`
	ShouldDisplayExecutiveSummarySection          bool               `json:"shouldDisplayExecutiveSummarySection"`
	ShouldDisplayLocationSection                  bool               `json:"shouldDisplayLocationSection"`
	ShouldDisplayNeighbourhoodSection             bool               `json:"shouldDisplayNeighbourhoodSection"`
	ShouldDisplayExteriorSection                  bool               `json:"shouldDisplayExteriorSection"`
	ShouldDisplayBuildingSection                  bool               `json:"shouldDisplayBuildingSection"`
	ShouldDisplayCommercialInteriorSection        bool               `json:"shouldDisplayCommercialInteriorSection"`
	ShouldDisplayCommercialInteriorDetailsSection bool               `json:"shouldDisplayCommercialInteriorDetailsSection"`
	ShouldDisplayResidentialInteriorSection       bool               `json:"shouldDisplayResidentialInteriorSection"`
	ShouldDisplayLegalSection                     bool               `json:"shouldDisplayLegalSection"`
	ShouldDisplayFinancialSection                 bool               `json:"shouldDisplayFinancialSection"`
	NeighbourhoodID                               primitive.ObjectID `json:"neighbourhoodId"`
}

type CreateEvaluationService interface {
	Execute(ctx context.Context, request *CreateEvaluationRequestDTO) (*CreateEvaluationResponseDTO, error)
}

type createEvaluationServiceImpl struct {
	config                  *config.Configuration
	logger                  *zap.Logger
	evaluationCreateUseCase uc_evaluation.EvaluationCreateUseCase
}

func NewCreateEvaluationService(
	config *config.Configuration,
	logger *zap.Logger,
	evaluationCreateUseCase uc_evaluation.EvaluationCreateUseCase,
) CreateEvaluationService {
	return &createEvaluationServiceImpl{
		config:                  config,
		logger:                  logger,
		evaluationCreateUseCase: evaluationCreateUseCase,
	}
}

func (svc *createEvaluationServiceImpl) Execute(ctx context.Context, req *CreateEvaluationRequestDTO) (*CreateEvaluationResponseDTO, error) {
	// Validate request
	if req == nil {
		return nil, httperror.NewForBadRequestWithSingleField("request", "Request is required")
	}

	errors := make(map[string]string)
	if req.PropertyID.IsZero() {
		errors["propertyId"] = "Property ID is required"
	}
	if req.ClientID.IsZero() {
		errors["clientId"] = "Client ID is required"
	}
	if req.PresenterID.IsZero() {
		errors["presenterId"] = "Presenter ID is required"
	}
	if req.OwnerID.IsZero() {
		errors["ownerId"] = "Owner ID is required"
	}
	// If neighbourhood section is displayed, neighbourhood ID is required
	if req.ShouldDisplayNeighbourhoodSection && req.NeighbourhoodID.IsZero() {
		errors["neighbourhoodId"] = "Neighbourhood ID is required when neighbourhood section is displayed"
	}

	if len(errors) > 0 {
		return nil, httperror.NewForBadRequest(&errors)
	}

	// Create evaluation domain object
	evaluation := &dom_evaluation.Evaluation{
		ID:                                     primitive.NewObjectID(),
		PropertyID:                             req.PropertyID,
		ClientID:                               req.ClientID,
		PresenterID:                            req.PresenterID,
		OwnerID:                                req.OwnerID,
		ShouldDisplayTitleSection:              req.ShouldDisplayTitleSection,
		ShouldDisplayExecutiveSummarySection:   req.ShouldDisplayExecutiveSummarySection,
		ShouldDisplayLocationSection:           req.ShouldDisplayLocationSection,
		ShouldDisplayNeighbourhoodSection:      req.ShouldDisplayNeighbourhoodSection,
		ShouldDisplayExteriorSection:           req.ShouldDisplayExteriorSection,
		ShouldDisplayBuildingSection:           req.ShouldDisplayBuildingSection,
		ShouldDisplayCommercialInteriorSection: req.ShouldDisplayCommercialInteriorSection,
		ShouldDisplayCommercialInteriorDetailsSection: req.ShouldDisplayCommercialInteriorDetailsSection,
		ShouldDisplayResidentialInteriorSection:       req.ShouldDisplayResidentialInteriorSection,
		ShouldDisplayLegalSection:                     req.ShouldDisplayLegalSection,
		ShouldDisplayFinancialSection:                 req.ShouldDisplayFinancialSection,
		NeighbourhoodID:                               req.NeighbourhoodID,
		Building:                                      dom_evaluation.Building{}, // Default empty building
		Legal:                                         dom_evaluation.Legal{},    // Default empty legal
	}

	// Create evaluation
	id, err := svc.evaluationCreateUseCase.Execute(ctx, evaluation)
	if err != nil {
		svc.logger.Error("Failed to create evaluation", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateEvaluationResponseDTO{
		ID:                                     id,
		PropertyID:                             evaluation.PropertyID,
		ClientID:                               evaluation.ClientID,
		PresenterID:                            evaluation.PresenterID,
		OwnerID:                                evaluation.OwnerID,
		ShouldDisplayTitleSection:              evaluation.ShouldDisplayTitleSection,
		ShouldDisplayExecutiveSummarySection:   evaluation.ShouldDisplayExecutiveSummarySection,
		ShouldDisplayLocationSection:           evaluation.ShouldDisplayLocationSection,
		ShouldDisplayNeighbourhoodSection:      evaluation.ShouldDisplayNeighbourhoodSection,
		ShouldDisplayExteriorSection:           evaluation.ShouldDisplayExteriorSection,
		ShouldDisplayBuildingSection:           evaluation.ShouldDisplayBuildingSection,
		ShouldDisplayCommercialInteriorSection: evaluation.ShouldDisplayCommercialInteriorSection,
		ShouldDisplayCommercialInteriorDetailsSection: evaluation.ShouldDisplayCommercialInteriorDetailsSection,
		ShouldDisplayResidentialInteriorSection:       evaluation.ShouldDisplayResidentialInteriorSection,
		ShouldDisplayLegalSection:                     evaluation.ShouldDisplayLegalSection,
		ShouldDisplayFinancialSection:                 evaluation.ShouldDisplayFinancialSection,
		NeighbourhoodID:                               evaluation.NeighbourhoodID,
	}

	return response, nil
}
