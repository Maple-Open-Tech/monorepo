// cloud/backend/internal/ipe/service/evaluation/get.go
package evaluation

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
	uc_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/evaluation"
)

type EvaluationResponseDTO struct {
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
	Building                                      BuildingDTO        `json:"building"`
	Legal                                         LegalDTO           `json:"legal"`
	NeighbourhoodID                               primitive.ObjectID `json:"neighbourhoodId"`
	PropertyPhotos                                []PropertyPhotoDTO `json:"propertyPhotos,omitempty"`
}

type BuildingDTO struct {
	Basement                              string  `json:"basement"`
	BuildingDesign                        string  `json:"buildingDesign"`
	BuildingStyle                         string  `json:"buildingStyle"`
	BuildingType                          string  `json:"buildingType"`
	Ceiling                               string  `json:"ceiling"`
	CeilingHeightInFeet                   float64 `json:"ceilingHeightInFeet"`
	CommercialAccess                      string  `json:"commercialAccess"`
	CommercialCondition                   string  `json:"commercialCondition"`
	CommercialDescription                 string  `json:"commercialDescription"`
	CommercialGrossAreaInSquareFeet       float64 `json:"commercialGrossAreaInSquareFeet"`
	CommercialNetRentableAreaInSquareFeet float64 `json:"commercialNetRentableAreaInSquareFeet"`
	CommercialType                        string  `json:"commercialType"`
	CommercialUnits                       string  `json:"commercialUnits"`
	CoolingSystem                         string  `json:"coolingSystem"`
	DeferredMaintenance                   string  `json:"deferredMaintenance"`
	ElectricalSystem                      string  `json:"electricalSystem"`
	ExpectedUsefulLife                    float64 `json:"expectedUsefulLife"`
	ExteriorDoorMaterial                  string  `json:"exteriorDoorMaterial"`
	ExteriorWallMaterial                  string  `json:"exteriorWallMaterial"`
	FireSystem                            string  `json:"fireSystem"`
	FloorCover                            string  `json:"floorCover"`
	Footing                               string  `json:"footing"`
	FoundationWall                        string  `json:"foundationWall"`
	Framing                               string  `json:"framing"`
	FunctionalUtility                     string  `json:"functionalUtility"`
	GrossBuildingAreaInSquareFeet         float64 `json:"grossBuildingAreaInSquareFeet"`
	GrossLandAreaInSquareFeet             float64 `json:"grossLandAreaInSquareFeet"`
	HeatingSystem                         string  `json:"heatingSystem"`
	OverallExteriorCondition              string  `json:"overallExteriorCondition"`
	PartitionWall                         string  `json:"partitionWall"`
	Plumbing                              string  `json:"plumbing"`
	RoofConstruction                      string  `json:"roofConstruction"`
	RoofStyle                             string  `json:"roofStyle"`
	SafetySystem                          string  `json:"safetySystem"`
	SiteCoverageRatio                     float64 `json:"siteCoverageRatio"`
	Stories                               float64 `json:"stories"`
	TotalFullBathRooms                    int     `json:"totalFullBathRooms"`
	TotalFullBedRooms                     int     `json:"totalFullBedRooms"`
	TotalHalfBathRooms                    int     `json:"totalHalfBathRooms"`
	TotalHalfBedRooms                     int     `json:"totalHalfBedRooms"`
	TotalNumberOfFamilyUnits              int     `json:"totalNumberOfFamilyUnits"`
	WindowType                            string  `json:"windowType"`
	YearBuilt                             int     `json:"yearBuilt"`
}

type LegalDTO struct {
	BuildingType               string  `json:"buildingType"`
	Designation                string  `json:"designation"`
	Fencing                    string  `json:"fencing"`
	FrontageInFeet             float64 `json:"frontageInFeet"`
	HasSoldWithinPastFiveYears bool    `json:"hasSoldWithinPastFiveYears"`
	Landscaping                string  `json:"landscaping"`
	LegalDescription           string  `json:"legalDescription"`
	Lighting                   string  `json:"lighting"`
	ParkingSpaces              int     `json:"parkingSpaces"`
	PermittedUses              string  `json:"permittedUses"`
	PhaseInAssessedValue       string  `json:"phaseInAssessedValue"`
	RollNumber                 float64 `json:"rollNumber"`
	ShapeOfLandParcel          string  `json:"shapeOfLandParcel"`
	TaxYear                    int     `json:"taxYear"`
	Topography                 string  `json:"topography"`
	TotalPropertyAreaInAcres   float64 `json:"totalPropertyAreaInAcres"`
	TotalTaxes                 float64 `json:"totalTaxes"`
	ZoneCode                   string  `json:"zoneCode"`
}

type PropertyPhotoDTO struct {
	ID             primitive.ObjectID `json:"id"`
	EvaluationID   primitive.ObjectID `json:"evaluationId"`
	PhotoCategory  string             `json:"photoCategory"`
	PhotoComment   string             `json:"photoComment"`
	HasPhotoData   bool               `json:"hasPhotoData"`
	PhotoName      string             `json:"photoName"`
	PhotoTimestamp int64              `json:"photoTimestamp"`
	PhotoUniqueID  string             `json:"photoUniqueId"`
}

type GetEvaluationService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*EvaluationResponseDTO, error)
	ExecuteByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*EvaluationResponseDTO, error)
}

type getEvaluationServiceImpl struct {
	config                           *config.Configuration
	logger                           *zap.Logger
	evaluationGetByIDUseCase         uc_evaluation.EvaluationGetByIDUseCase
	evaluationGetByPropertyIDUseCase uc_evaluation.EvaluationGetByPropertyIDUseCase
	findPhotosByEvaluationIDUseCase  uc_evaluation.FindPhotosByEvaluationIDUseCase
}

func NewGetEvaluationService(
	config *config.Configuration,
	logger *zap.Logger,
	evaluationGetByIDUseCase uc_evaluation.EvaluationGetByIDUseCase,
	evaluationGetByPropertyIDUseCase uc_evaluation.EvaluationGetByPropertyIDUseCase,
	findPhotosByEvaluationIDUseCase uc_evaluation.FindPhotosByEvaluationIDUseCase,
) GetEvaluationService {
	return &getEvaluationServiceImpl{
		config:                           config,
		logger:                           logger,
		evaluationGetByIDUseCase:         evaluationGetByIDUseCase,
		evaluationGetByPropertyIDUseCase: evaluationGetByPropertyIDUseCase,
		findPhotosByEvaluationIDUseCase:  findPhotosByEvaluationIDUseCase,
	}
}

func (svc *getEvaluationServiceImpl) Execute(ctx context.Context, id primitive.ObjectID) (*EvaluationResponseDTO, error) {
	// Get evaluation
	evaluation, err := svc.evaluationGetByIDUseCase.Execute(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Evaluation not found", zap.String("id", id.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get evaluation", zap.Error(err))
		return nil, err
	}
	if evaluation == nil {
		svc.logger.Warn("Evaluation not found", zap.String("id", id.Hex()))
		return nil, errors.New("evaluation not found")
	}

	// Get photos
	photos, err := svc.findPhotosByEvaluationIDUseCase.Execute(ctx, id)
	if err != nil {
		svc.logger.Error("Failed to get evaluation photos", zap.Error(err))
		return nil, err
	}

	// Map photos
	propertyPhotos := make([]PropertyPhotoDTO, len(photos))
	for i, photo := range photos {
		propertyPhotos[i] = PropertyPhotoDTO{
			ID:             photo.ID,
			EvaluationID:   photo.EvaluationID,
			PhotoCategory:  photo.PhotoCategory,
			PhotoComment:   photo.PhotoComment,
			HasPhotoData:   len(photo.PhotoData) > 0,
			PhotoName:      photo.PhotoName,
			PhotoTimestamp: photo.PhotoTimestamp,
			PhotoUniqueID:  photo.PhotoUniqueID,
		}
	}

	// Create response
	response := &EvaluationResponseDTO{
		ID:                                     evaluation.ID,
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
		Building:                                      mapBuildingToDTO(evaluation.Building),
		Legal:                                         mapLegalToDTO(evaluation.Legal),
		NeighbourhoodID:                               evaluation.NeighbourhoodID,
		PropertyPhotos:                                propertyPhotos,
	}

	return response, nil
}

func (svc *getEvaluationServiceImpl) ExecuteByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*EvaluationResponseDTO, error) {
	// Get evaluation
	evaluation, err := svc.evaluationGetByPropertyIDUseCase.Execute(ctx, propertyID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Evaluation not found for property", zap.String("propertyId", propertyID.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get evaluation by property ID", zap.Error(err))
		return nil, err
	}
	if evaluation == nil {
		svc.logger.Warn("Evaluation not found for property", zap.String("propertyId", propertyID.Hex()))
		return nil, errors.New("evaluation not found")
	}

	return svc.Execute(ctx, evaluation.ID)
}

func mapBuildingToDTO(building dom_evaluation.Building) BuildingDTO {
	return BuildingDTO{
		Basement:                              building.Basement,
		BuildingDesign:                        building.BuildingDesign,
		BuildingStyle:                         building.BuildingStyle,
		BuildingType:                          building.BuildingType,
		Ceiling:                               building.Ceiling,
		CeilingHeightInFeet:                   building.CeilingHeightInFeet,
		CommercialAccess:                      building.CommercialAccess,
		CommercialCondition:                   building.CommercialCondition,
		CommercialDescription:                 building.CommercialDescription,
		CommercialGrossAreaInSquareFeet:       building.CommercialGrossAreaInSquareFeet,
		CommercialNetRentableAreaInSquareFeet: building.CommercialNetRentableAreaInSquareFeet,
		CommercialType:                        building.CommercialType,
		CommercialUnits:                       building.CommercialUnits,
		CoolingSystem:                         building.CoolingSystem,
		DeferredMaintenance:                   building.DeferredMaintenance,
		ElectricalSystem:                      building.ElectricalSystem,
		ExpectedUsefulLife:                    building.ExpectedUsefulLife,
		ExteriorDoorMaterial:                  building.ExteriorDoorMaterial,
		ExteriorWallMaterial:                  building.ExteriorWallMaterial,
		FireSystem:                            building.FireSystem,
		FloorCover:                            building.FloorCover,
		Footing:                               building.Footing,
		FoundationWall:                        building.FoundationWall,
		Framing:                               building.Framing,
		FunctionalUtility:                     building.FunctionalUtility,
		GrossBuildingAreaInSquareFeet:         building.GrossBuildingAreaInSquareFeet,
		GrossLandAreaInSquareFeet:             building.GrossLandAreaInSquareFeet,
		HeatingSystem:                         building.HeatingSystem,
		OverallExteriorCondition:              building.OverallExteriorCondition,
		PartitionWall:                         building.PartitionWall,
		Plumbing:                              building.Plumbing,
		RoofConstruction:                      building.RoofConstruction,
		RoofStyle:                             building.RoofStyle,
		SafetySystem:                          building.SafetySystem,
		SiteCoverageRatio:                     building.SiteCoverageRatio,
		Stories:                               building.Stories,
		TotalFullBathRooms:                    building.TotalFullBathRooms,
		TotalFullBedRooms:                     building.TotalFullBedRooms,
		TotalHalfBathRooms:                    building.TotalHalfBathRooms,
		TotalHalfBedRooms:                     building.TotalHalfBedRooms,
		TotalNumberOfFamilyUnits:              building.TotalNumberOfFamilyUnits,
		WindowType:                            building.WindowType,
		YearBuilt:                             building.YearBuilt,
	}
}

func mapLegalToDTO(legal dom_evaluation.Legal) LegalDTO {
	return LegalDTO{
		BuildingType:               legal.BuildingType,
		Designation:                legal.Designation,
		Fencing:                    legal.Fencing,
		FrontageInFeet:             legal.FrontageInFeet,
		HasSoldWithinPastFiveYears: legal.HasSoldWithinPastFiveYears,
		Landscaping:                legal.Landscaping,
		LegalDescription:           legal.LegalDescription,
		Lighting:                   legal.Lighting,
		ParkingSpaces:              legal.ParkingSpaces,
		PermittedUses:              legal.PermittedUses,
		PhaseInAssessedValue:       legal.PhaseInAssessedValue.String(),
		RollNumber:                 legal.RollNumber,
		ShapeOfLandParcel:          legal.ShapeOfLandParcel,
		TaxYear:                    legal.TaxYear,
		Topography:                 legal.Topography,
		TotalPropertyAreaInAcres:   legal.TotalPropertyAreaInAcres,
		TotalTaxes:                 legal.TotalTaxes,
		ZoneCode:                   legal.ZoneCode,
	}
}
