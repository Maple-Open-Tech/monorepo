package incomeproperty

import (
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IncomeProperty represents a comprehensive real estate income property
type IncomeProperty struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Address            string             `bson:"address" json:"address"`
	City               string             `bson:"city" json:"city"`
	Province           string             `bson:"province" json:"province"`
	Country            string             `bson:"country" json:"country"`
	PropertyCode       string             `bson:"propertyCode" json:"propertyCode"`
	RecordName         string             `bson:"recordName" json:"recordName"`
	RecordCreationDate time.Time          `bson:"recordCreationDate" json:"recordCreationDate"`
	MainPhotoThumbnail []byte             `bson:"mainPhotoThumbnail,omitempty" json:"mainPhotoThumbnail,omitempty"`

	// Embedded evaluation
	Evaluation *Evaluation `bson:"evaluation,omitempty" json:"evaluation,omitempty"`

	// Embedded financial analysis
	FinancialAnalysis *FinancialAnalysis `bson:"financialAnalysis,omitempty" json:"financialAnalysis,omitempty"`

	// Embedded persons
	Client    *Client    `bson:"client,omitempty" json:"client,omitempty"`
	Presenter *Presenter `bson:"presenter,omitempty" json:"presenter,omitempty"`
	Owner     *Owner     `bson:"owner,omitempty" json:"owner,omitempty"`

	// Embedded comparisons
	Comparisons []*Compare `bson:"comparisons,omitempty" json:"comparisons,omitempty"`
}

// Evaluation represents property evaluation details
type Evaluation struct {
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
	Building      Building       `bson:"building" json:"building"`
	Legal         Legal          `bson:"legal" json:"legal"`
	Neighbourhood *Neighbourhood `bson:"neighbourhood,omitempty" json:"neighbourhood,omitempty"`

	// Embedded property photos
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
	RecordName              string `bson:"recordName" json:"recordName"`
	RecordUniqueID          string `bson:"recordUniqueId" json:"recordUniqueId"`
	Appeal                  string `bson:"appeal" json:"appeal"`
	City                    string `bson:"city" json:"city"`
	Province                string `bson:"province" json:"province"`
	Country                 string `bson:"country" json:"country"`
	Comment                 string `bson:"comment" json:"comment"`
	DevelopmentTrend        string `bson:"developmentTrend" json:"developmentTrend"`
	DominantLandUse         string `bson:"dominantLandUse" json:"dominantLandUse"`
	AdditionalLandUse       string `bson:"additionalLandUse" json:"additionalLandUse"`
	EstablishedYear         int    `bson:"establishedYear" json:"establishedYear"`
	GeneralValueTrend       string `bson:"generalValueTrend" json:"generalValueTrend"`
	HasCurbsAndGutters      bool   `bson:"hasCurbsAndGutters" json:"hasCurbsAndGutters"`
	HasPublicTransportation bool   `bson:"hasPublicTransportation" json:"hasPublicTransportation"`
	HasSideWalks            bool   `bson:"hasSideWalks" json:"hasSideWalks"`
	LocationWithinCity      string `bson:"locationWithinCity" json:"locationWithinCity"`
	PopulationTrend         string `bson:"populationTrend" json:"populationTrend"`
	StandardMapPhotoData    []byte `bson:"standardMapPhotoData,omitempty" json:"standardMapPhotoData,omitempty"`
}

// PropertyPhoto represents property photos
type PropertyPhoto struct {
	PhotoCategory  string `bson:"photoCategory" json:"photoCategory"`
	PhotoComment   string `bson:"photoComment" json:"photoComment"`
	PhotoData      []byte `bson:"photoData,omitempty" json:"photoData,omitempty"`
	PhotoName      string `bson:"photoName" json:"photoName"`
	PhotoTimestamp int64  `bson:"photoTimestamp" json:"photoTimestamp"`
	PhotoUniqueID  string `bson:"photoUniqueId" json:"photoUniqueId"`
}

// FinancialAnalysis represents the financial analysis for a property
type FinancialAnalysis struct {
	PurchasePrice             decimal.Decimal `bson:"purchasePrice" json:"purchasePrice"`
	AnnualGrossIncome         decimal.Decimal `bson:"annualGrossIncome" json:"annualGrossIncome"`
	MonthlyGrossIncome        decimal.Decimal `bson:"monthlyGrossIncome" json:"monthlyGrossIncome"`
	AnnualExpense             decimal.Decimal `bson:"annualExpense" json:"annualExpense"`
	MonthlyExpense            decimal.Decimal `bson:"monthlyExpense" json:"monthlyExpense"`
	AnnualNetIncome           decimal.Decimal `bson:"annualNetIncome" json:"annualNetIncome"`
	MonthlyNetIncome          decimal.Decimal `bson:"monthlyNetIncome" json:"monthlyNetIncome"`
	AnnualCashFlow            decimal.Decimal `bson:"annualCashFlow" json:"annualCashFlow"`
	MonthlyCashFlow           decimal.Decimal `bson:"monthlyCashFlow" json:"monthlyCashFlow"`
	CapRateWithMortgage       decimal.Decimal `bson:"capRateWithMortgage" json:"capRateWithMortgage"`
	CapRateWithoutMortgage    decimal.Decimal `bson:"capRateWithoutMortgage" json:"capRateWithoutMortgage"`
	AnnualRentalIncome        decimal.Decimal `bson:"annualRentalIncome" json:"annualRentalIncome"`
	MonthlyRentalIncome       decimal.Decimal `bson:"monthlyRentalIncome" json:"monthlyRentalIncome"`
	AnnualFacilityIncome      decimal.Decimal `bson:"annualFacilityIncome" json:"annualFacilityIncome"`
	MonthlyFacilityIncome     decimal.Decimal `bson:"monthlyFacilityIncome" json:"monthlyFacilityIncome"`
	BuyingFeeRate             decimal.Decimal `bson:"buyingFeeRate" json:"buyingFeeRate"`
	SellingFeeRate            decimal.Decimal `bson:"sellingFeeRate" json:"sellingFeeRate"`
	CapitalImprovementsAmount decimal.Decimal `bson:"capitalImprovementsAmount" json:"capitalImprovementsAmount"`
	PurchaseFeesAmount        decimal.Decimal `bson:"purchaseFeesAmount" json:"purchaseFeesAmount"`
	InitialInvestmentAmount   decimal.Decimal `bson:"initialInvestmentAmount" json:"initialInvestmentAmount"`
	InflationRate             decimal.Decimal `bson:"inflationRate" json:"inflationRate"`

	// Embedded related collections
	RentalIncomes       []RentalIncome       `bson:"rentalIncomes,omitempty" json:"rentalIncomes,omitempty"`
	CommercialIncomes   []CommercialIncome   `bson:"commercialIncomes,omitempty" json:"commercialIncomes,omitempty"`
	FacilityIncomes     []FacilityIncome     `bson:"facilityIncomes,omitempty" json:"facilityIncomes,omitempty"`
	Expenses            []Expense            `bson:"expenses,omitempty" json:"expenses,omitempty"`
	PurchaseFees        []PurchaseFee        `bson:"purchaseFees,omitempty" json:"purchaseFees,omitempty"`
	CapitalImprovements []CapitalImprovement `bson:"capitalImprovements,omitempty" json:"capitalImprovements,omitempty"`
	AnnualProjections   []AnnualProjection   `bson:"annualProjections,omitempty" json:"annualProjections,omitempty"`
	Mortgage            *Mortgage            `bson:"mortgage,omitempty" json:"mortgage,omitempty"`
}

// RentalIncome represents income from rental units
type RentalIncome struct {
	NameText             string          `bson:"nameText" json:"nameText"`
	MonthlyAmount        decimal.Decimal `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount         decimal.Decimal `bson:"annualAmount" json:"annualAmount"`
	MonthlyAmountPerUnit decimal.Decimal `bson:"monthlyAmountPerUnit" json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  decimal.Decimal `bson:"annualAmountPerUnit" json:"annualAmountPerUnit"`
	NumberOfUnits        decimal.Decimal `bson:"numberOfUnits" json:"numberOfUnits"`
	Frequency            decimal.Decimal `bson:"frequency" json:"frequency"`
	TypeID               int             `bson:"typeId" json:"typeId"`
}

// CommercialIncome represents income from commercial spaces
type CommercialIncome struct {
	NameText             string          `bson:"nameText" json:"nameText"`
	MonthlyAmount        decimal.Decimal `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount         decimal.Decimal `bson:"annualAmount" json:"annualAmount"`
	MonthlyAmountPerUnit decimal.Decimal `bson:"monthlyAmountPerUnit" json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  decimal.Decimal `bson:"annualAmountPerUnit" json:"annualAmountPerUnit"`
	AreaInSquareFeet     decimal.Decimal `bson:"areaInSquareFeet" json:"areaInSquareFeet"`
	UnitType             string          `bson:"unitType" json:"unitType"`
	UnitValue            decimal.Decimal `bson:"unitValue" json:"unitValue"`
	Frequency            decimal.Decimal `bson:"frequency" json:"frequency"`
	TypeID               int             `bson:"typeId" json:"typeId"`
}

// FacilityIncome represents income from property facilities
type FacilityIncome struct {
	NameText      string          `bson:"nameText" json:"nameText"`
	MonthlyAmount decimal.Decimal `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount  decimal.Decimal `bson:"annualAmount" json:"annualAmount"`
	Frequency     decimal.Decimal `bson:"frequency" json:"frequency"`
	TypeID        int             `bson:"typeId" json:"typeId"`
}

// Expense represents property expenses
type Expense struct {
	NameText      string          `bson:"nameText" json:"nameText"`
	MonthlyAmount decimal.Decimal `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount  decimal.Decimal `bson:"annualAmount" json:"annualAmount"`
	Frequency     decimal.Decimal `bson:"frequency" json:"frequency"`
	Percent       decimal.Decimal `bson:"percent" json:"percent"`
	TypeID        int             `bson:"typeId" json:"typeId"`
}

// PurchaseFee represents fees associated with property purchase
type PurchaseFee struct {
	NameText string          `bson:"nameText" json:"nameText"`
	Amount   decimal.Decimal `bson:"amount" json:"amount"`
	TypeID   int             `bson:"typeId" json:"typeId"`
}

// CapitalImprovement represents capital improvements to property
type CapitalImprovement struct {
	NameText string          `bson:"nameText" json:"nameText"`
	Amount   decimal.Decimal `bson:"amount" json:"amount"`
	TypeID   int             `bson:"typeId" json:"typeId"`
}

// AnnualProjection represents annual financial projections
type AnnualProjection struct {
	Year                                decimal.Decimal `bson:"year" json:"year"`
	CashFlow                            decimal.Decimal `bson:"cashFlow" json:"cashFlow"`
	DebtRemaining                       decimal.Decimal `bson:"debtRemaining" json:"debtRemaining"`
	SalesPrice                          decimal.Decimal `bson:"salesPrice" json:"salesPrice"`
	LegalFees                           decimal.Decimal `bson:"legalFees" json:"legalFees"`
	ProceedsOfSale                      decimal.Decimal `bson:"proceedsOfSale" json:"proceedsOfSale"`
	TotalReturn                         decimal.Decimal `bson:"totalReturn" json:"totalReturn"`
	InitialInvestment                   decimal.Decimal `bson:"initialInvestment" json:"initialInvestment"`
	ReturnOnInvestmentRate              decimal.Decimal `bson:"returnOnInvestmentRate" json:"returnOnInvestmentRate"`
	ReturnOnInvestmentPercent           decimal.Decimal `bson:"returnOnInvestmentPercent" json:"returnOnInvestmentPercent"`
	AnnualizedReturnOnInvestmentRate    decimal.Decimal `bson:"annualizedReturnOnInvestmentRate" json:"annualizedReturnOnInvestmentRate"`
	AnnualizedReturnOnInvestmentPercent decimal.Decimal `bson:"annualizedReturnOnInvestmentPercent" json:"annualizedReturnOnInvestmentPercent"`
}

// Mortgage represents mortgage details for a property
type Mortgage struct {
	LoanAmount                         decimal.Decimal `bson:"loanAmount" json:"loanAmount"`
	LoanPurchaseAmount                 decimal.Decimal `bson:"loanPurchaseAmount" json:"loanPurchaseAmount"`
	DownPayment                        decimal.Decimal `bson:"downPayment" json:"downPayment"`
	AnnualInterestRate                 decimal.Decimal `bson:"annualInterestRate" json:"annualInterestRate"`
	AmortizationYear                   decimal.Decimal `bson:"amortizationYear" json:"amortizationYear"`
	PaymentFrequency                   decimal.Decimal `bson:"paymentFrequency" json:"paymentFrequency"`
	CompoundingPeriod                  decimal.Decimal `bson:"compoundingPeriod" json:"compoundingPeriod"`
	FirstPaymentDate                   time.Time       `bson:"firstPaymentDate" json:"firstPaymentDate"`
	Insurance                          string          `bson:"insurance" json:"insurance"`
	InsuranceAmount                    decimal.Decimal `bson:"insuranceAmount" json:"insuranceAmount"`
	MortgagePaymentPerPaymentFrequency decimal.Decimal `bson:"mortgagePaymentPerPaymentFrequency" json:"mortgagePaymentPerPaymentFrequency"`
	InterestRatePerPaymentFrequency    decimal.Decimal `bson:"interestRatePerPaymentFrequency" json:"interestRatePerPaymentFrequency"`
	TotalNumberOfPaymentsPerFrequency  decimal.Decimal `bson:"totalNumberOfPaymentsPerFrequency" json:"totalNumberOfPaymentsPerFrequency"`
	PercentFinanced                    decimal.Decimal `bson:"percentFinanced" json:"percentFinanced"`

	// Embedded payment schedule
	MortgagePaymentSchedule []MortgageInterval `bson:"mortgagePaymentSchedule,omitempty" json:"mortgagePaymentSchedule,omitempty"`
}

// MortgageInterval represents a payment in the mortgage schedule
type MortgageInterval struct {
	Interval            decimal.Decimal `bson:"interval" json:"interval"`
	PaymentDate         time.Time       `bson:"paymentDate" json:"paymentDate"`
	PaymentAmount       decimal.Decimal `bson:"paymentAmount" json:"paymentAmount"`
	PrincipleAmount     decimal.Decimal `bson:"principleAmount" json:"principleAmount"`
	InterestAmount      decimal.Decimal `bson:"interestAmount" json:"interestAmount"`
	LoanBalance         decimal.Decimal `bson:"loanBalance" json:"loanBalance"`
	Year                decimal.Decimal `bson:"year" json:"year"`
	TotalPaidToBank     decimal.Decimal `bson:"totalPaidToBank" json:"totalPaidToBank"`
	TotalPaidToInterest decimal.Decimal `bson:"totalPaidToInterest" json:"totalPaidToInterest"`
}

// Common fields for all person types
type basePerson struct {
	PersonName      string `bson:"personName" json:"personName"`
	Address         string `bson:"address" json:"address"`
	City            string `bson:"city" json:"city"`
	Province        string `bson:"province" json:"province"`
	Country         string `bson:"country" json:"country"`
	PostalCode      string `bson:"postalCode" json:"postalCode"`
	Email           string `bson:"email" json:"email"`
	OfficeTelNumber string `bson:"officeTelNumber" json:"officeTelNumber"`
	MobileTelNumber string `bson:"mobileTelNumber" json:"mobileTelNumber"`
	FaxTelNumber    string `bson:"faxTelNumber" json:"faxTelNumber"`
	Website         string `bson:"website" json:"website"`
	RecordUniqueID  string `bson:"recordUniqueId" json:"recordUniqueId"`
	LogoPhotoData   []byte `bson:"logoPhotoData,omitempty" json:"logoPhotoData,omitempty"`
}

// Client represents a client
type Client struct {
	basePerson
}

// Presenter represents a presenter
type Presenter struct {
	basePerson
	RecordUnquieID string `bson:"recordUnquieId" json:"recordUnquieId"`
}

// Owner represents a property owner
type Owner struct {
	basePerson
}

// Compare represents a property comparison
type Compare struct {
	RecordUniqueID           string               `bson:"recordUniqueId" json:"recordUniqueId"`
	SelectedIncomeProperties []primitive.ObjectID `bson:"selectedIncomeProperties" json:"selectedIncomeProperties"`
}
