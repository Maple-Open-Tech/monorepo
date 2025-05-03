package evaluation

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Evaluation represents a property evaluation
type Evaluation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PropertyID  primitive.ObjectID `bson:"propertyId" json:"propertyId"`
	ClientID    primitive.ObjectID `bson:"clientId" json:"clientId"`
	PresenterID primitive.ObjectID `bson:"presenterId" json:"presenterId"`
	OwnerID     primitive.ObjectID `bson:"ownerId" json:"ownerId"`

	// Display flags for sections
	ShouldDisplayTitleSection                     bool `bson:"shouldDisplayTitleSection" json:"shouldDisplayTitleSection"`
	ShouldDisplayExecutiveSummarySection          bool `bson:"shouldDisplayExecutiveSummarySection" json:"shouldDisplayExecutiveSummarySection"`
	ShouldDisplayLocationSection                  bool `bson:"shouldDisplayLocationSection" json:"shouldDisplayLocationSection"`
	ShouldDisplayNeighbourhoodSection             bool `bson:"shouldDisplayNeighbourhoodSection" json:"shouldDisplayNeighbourhoodSection"`
	ShouldDisplayExteriorSection                  bool `bson:"shouldDisplayExteriorSection" json:"shouldDisplayExteriorSection"`
	ShouldDisplayBuildingSection                  bool `bson:"shouldDisplayBuildingSection" json:"shouldDisplayBuildingSection"`
	ShouldDisplayCommercialInteriorSection        bool `bson:"shouldDisplayCommercialInteriorSection" json:"shouldDisplayCommercialInteriorSection"`
	ShouldDisplayCommercialInteriorDetailsSection bool `bson:"shouldDisplayCommercialInteriorDetailsSection" json:"shouldDisplayCommercialInteriorDetailsSection"`
	ShouldDisplayResidentialInteriorSection       bool `bson:"shouldDisplayResidentialInteriorSection" json:"shouldDisplayResidentialInteriorSection"`
	ShouldDisplayLegalSection                     bool `bson:"shouldDisplayLegalSection" json:"shouldDisplayLegalSection"`
	ShouldDisplayFinancialSection                 bool `bson:"shouldDisplayFinancialSection" json:"shouldDisplayFinancialSection"`

	// Embedded related entities
	Building        Building           `bson:"building" json:"building"`
	Legal           Legal              `bson:"legal" json:"legal"`
	NeighbourhoodID primitive.ObjectID `bson:"neighbourhoodId" json:"neighbourhoodId"`

	// Reference to other collections
	PropertyPhotos []PropertyPhoto `bson:"propertyPhotos,omitempty" json:"propertyPhotos,omitempty"`
}

// Building represents building details
type Building struct {
	Basement                              string  `bson:"basement" json:"basement"`
	BuildingDesign                        string  `bson:"buildingDesign" json:"buildingDesign"`
	BuildingStyle                         string  `bson:"buildingStyle" json:"buildingStyle"`
	BuildingType                          string  `bson:"buildingType" json:"buildingType"`
	Ceiling                               string  `bson:"ceiling" json:"ceiling"`
	CeilingHeightInFeet                   float64 `bson:"ceilingHeightInFeet" json:"ceilingHeightInFeet"`
	CommercialAccess                      string  `bson:"commercialAccess" json:"commercialAccess"`
	CommercialCondition                   string  `bson:"commercialCondition" json:"commercialCondition"`
	CommercialDescription                 string  `bson:"commercialDescription" json:"commercialDescription"`
	CommercialGrossAreaInSquareFeet       float64 `bson:"commercialGrossAreaInSquareFeet" json:"commercialGrossAreaInSquareFeet"`
	CommercialNetRentableAreaInSquareFeet float64 `bson:"commercialNetRentableAreaInSquareFeet" json:"commercialNetRentableAreaInSquareFeet"`
	CommercialType                        string  `bson:"commercialType" json:"commercialType"`
	CommercialUnits                       string  `bson:"commercialUnits" json:"commercialUnits"`
	CoolingSystem                         string  `bson:"coolingSystem" json:"coolingSystem"`
	DeferredMaintenance                   string  `bson:"deferredMaintenance" json:"deferredMaintenance"`
	ElectricalSystem                      string  `bson:"electricalSystem" json:"electricalSystem"`
	ExpectedUsefulLife                    float64 `bson:"expectedUsefulLife" json:"expectedUsefulLife"`
	ExteriorDoorMaterial                  string  `bson:"exteriorDoorMaterial" json:"exteriorDoorMaterial"`
	ExteriorWallMaterial                  string  `bson:"exteriorWallMaterial" json:"exteriorWallMaterial"`
	FireSystem                            string  `bson:"fireSystem" json:"fireSystem"`
	FloorCover                            string  `bson:"floorCover" json:"floorCover"`
	Footing                               string  `bson:"footing" json:"footing"`
	FoundationWall                        string  `bson:"foundationWall" json:"foundationWall"`
	Framing                               string  `bson:"framing" json:"framing"`
	FunctionalUtility                     string  `bson:"functionalUtility" json:"functionalUtility"`
	GrossBuildingAreaInSquareFeet         float64 `bson:"grossBuildingAreaInSquareFeet" json:"grossBuildingAreaInSquareFeet"`
	GrossLandAreaInSquareFeet             float64 `bson:"grossLandAreaInSquareFeet" json:"grossLandAreaInSquareFeet"`
	HeatingSystem                         string  `bson:"heatingSystem" json:"heatingSystem"`
	OverallExteriorCondition              string  `bson:"overallExteriorCondition" json:"overallExteriorCondition"`
	PartitionWall                         string  `bson:"partitionWall" json:"partitionWall"`
	Plumbing                              string  `bson:"plumbing" json:"plumbing"`
	RoofConstruction                      string  `bson:"roofConstruction" json:"roofConstruction"`
	RoofStyle                             string  `bson:"roofStyle" json:"roofStyle"`
	SafetySystem                          string  `bson:"safetySystem" json:"safetySystem"`
	SiteCoverageRatio                     float64 `bson:"siteCoverageRatio" json:"siteCoverageRatio"`
	Stories                               float64 `bson:"stories" json:"stories"`
	TotalFullBathRooms                    int     `bson:"totalFullBathRooms" json:"totalFullBathRooms"`
	TotalFullBedRooms                     int     `bson:"totalFullBedRooms" json:"totalFullBedRooms"`
	TotalHalfBathRooms                    int     `bson:"totalHalfBathRooms" json:"totalHalfBathRooms"`
	TotalHalfBedRooms                     int     `bson:"totalHalfBedRooms" json:"totalHalfBedRooms"`
	TotalNumberOfFamilyUnits              int     `bson:"totalNumberOfFamilyUnits" json:"totalNumberOfFamilyUnits"`
	WindowType                            string  `bson:"windowType" json:"windowType"`
	YearBuilt                             int     `bson:"yearBuilt" json:"yearBuilt"`
}

// Legal represents legal details of a property
type Legal struct {
	BuildingType               string          `bson:"buildingType" json:"buildingType"`
	Designation                string          `bson:"designation" json:"designation"`
	Fencing                    string          `bson:"fencing" json:"fencing"`
	FrontageInFeet             float64         `bson:"frontageInFeet" json:"frontageInFeet"`
	HasSoldWithinPastFiveYears bool            `bson:"hasSoldWithinPastFiveYears" json:"hasSoldWithinPastFiveYears"`
	Landscaping                string          `bson:"landscaping" json:"landscaping"`
	LegalDescription           string          `bson:"legalDescription" json:"legalDescription"`
	Lighting                   string          `bson:"lighting" json:"lighting"`
	ParkingSpaces              int             `bson:"parkingSpaces" json:"parkingSpaces"`
	PermittedUses              string          `bson:"permittedUses" json:"permittedUses"`
	PhaseInAssessedValue       decimal.Decimal `bson:"phaseInAssessedValue" json:"phaseInAssessedValue"`
	RollNumber                 float64         `bson:"rollNumber" json:"rollNumber"`
	ShapeOfLandParcel          string          `bson:"shapeOfLandParcel" json:"shapeOfLandParcel"`
	TaxYear                    int             `bson:"taxYear" json:"taxYear"`
	Topography                 string          `bson:"topography" json:"topography"`
	TotalPropertyAreaInAcres   float64         `bson:"totalPropertyAreaInAcres" json:"totalPropertyAreaInAcres"`
	TotalTaxes                 float64         `bson:"totalTaxes" json:"totalTaxes"`
	ZoneCode                   string          `bson:"zoneCode" json:"zoneCode"`
}

// Neighbourhood represents neighbourhood details
type Neighbourhood struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RecordName              string             `bson:"recordName" json:"recordName"`
	RecordUniqueID          string             `bson:"recordUniqueId" json:"recordUniqueId"`
	Appeal                  string             `bson:"appeal" json:"appeal"`
	City                    string             `bson:"city" json:"city"`
	Province                string             `bson:"province" json:"province"`
	Country                 string             `bson:"country" json:"country"`
	Comment                 string             `bson:"comment" json:"comment"`
	DevelopmentTrend        string             `bson:"developmentTrend" json:"developmentTrend"`
	DominantLandUse         string             `bson:"dominantLandUse" json:"dominantLandUse"`
	AdditionalLandUse       string             `bson:"additionalLandUse" json:"additionalLandUse"`
	EstablishedYear         int                `bson:"establishedYear" json:"establishedYear"`
	GeneralValueTrend       string             `bson:"generalValueTrend" json:"generalValueTrend"`
	HasCurbsAndGutters      bool               `bson:"hasCurbsAndGutters" json:"hasCurbsAndGutters"`
	HasPublicTransportation bool               `bson:"hasPublicTransportation" json:"hasPublicTransportation"`
	HasSideWalks            bool               `bson:"hasSideWalks" json:"hasSideWalks"`
	LocationWithinCity      string             `bson:"locationWithinCity" json:"locationWithinCity"`
	PopulationTrend         string             `bson:"populationTrend" json:"populationTrend"`
	StandardMapPhotoData    []byte             `bson:"standardMapPhotoData,omitempty" json:"standardMapPhotoData,omitempty"`
}

// PropertyPhoto represents property photos
type PropertyPhoto struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EvaluationID   primitive.ObjectID `bson:"evaluationId" json:"evaluationId"`
	PhotoCategory  string             `bson:"photoCategory" json:"photoCategory"`
	PhotoComment   string             `bson:"photoComment" json:"photoComment"`
	PhotoData      []byte             `bson:"photoData,omitempty" json:"photoData,omitempty"`
	PhotoName      string             `bson:"photoName" json:"photoName"`
	PhotoTimestamp int64              `bson:"photoTimestamp" json:"photoTimestamp"`
	PhotoUniqueID  string             `bson:"photoUniqueId" json:"photoUniqueId"`
}
