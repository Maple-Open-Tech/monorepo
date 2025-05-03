package financialanalysis

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FinancialAnalysis represents the financial analysis for a property
type FinancialAnalysis struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PropertyID                primitive.ObjectID `bson:"propertyId" json:"propertyId"`
	PurchasePrice             decimal.Decimal    `bson:"purchasePrice" json:"purchasePrice"`
	AnnualGrossIncome         decimal.Decimal    `bson:"annualGrossIncome" json:"annualGrossIncome"`
	MonthlyGrossIncome        decimal.Decimal    `bson:"monthlyGrossIncome" json:"monthlyGrossIncome"`
	AnnualExpense             decimal.Decimal    `bson:"annualExpense" json:"annualExpense"`
	MonthlyExpense            decimal.Decimal    `bson:"monthlyExpense" json:"monthlyExpense"`
	AnnualNetIncome           decimal.Decimal    `bson:"annualNetIncome" json:"annualNetIncome"`
	MonthlyNetIncome          decimal.Decimal    `bson:"monthlyNetIncome" json:"monthlyNetIncome"`
	AnnualCashFlow            decimal.Decimal    `bson:"annualCashFlow" json:"annualCashFlow"`
	MonthlyCashFlow           decimal.Decimal    `bson:"monthlyCashFlow" json:"monthlyCashFlow"`
	CapRateWithMortgage       decimal.Decimal    `bson:"capRateWithMortgage" json:"capRateWithMortgage"`
	CapRateWithoutMortgage    decimal.Decimal    `bson:"capRateWithoutMortgage" json:"capRateWithoutMortgage"`
	AnnualRentalIncome        decimal.Decimal    `bson:"annualRentalIncome" json:"annualRentalIncome"`
	MonthlyRentalIncome       decimal.Decimal    `bson:"monthlyRentalIncome" json:"monthlyRentalIncome"`
	AnnualFacilityIncome      decimal.Decimal    `bson:"annualFacilityIncome" json:"annualFacilityIncome"`
	MonthlyFacilityIncome     decimal.Decimal    `bson:"monthlyFacilityIncome" json:"monthlyFacilityIncome"`
	BuyingFeeRate             decimal.Decimal    `bson:"buyingFeeRate" json:"buyingFeeRate"`
	SellingFeeRate            decimal.Decimal    `bson:"sellingFeeRate" json:"sellingFeeRate"`
	CapitalImprovementsAmount decimal.Decimal    `bson:"capitalImprovementsAmount" json:"capitalImprovementsAmount"`
	PurchaseFeesAmount        decimal.Decimal    `bson:"purchaseFeesAmount" json:"purchaseFeesAmount"`
	InitialInvestmentAmount   decimal.Decimal    `bson:"initialInvestmentAmount" json:"initialInvestmentAmount"`
	InflationRate             decimal.Decimal    `bson:"inflationRate" json:"inflationRate"`

	// Embedded related collections
	RentalIncomes       []RentalIncome       `bson:"rentalIncomes,omitempty" json:"rentalIncomes,omitempty"`
	CommercialIncomes   []CommercialIncome   `bson:"commercialIncomes,omitempty" json:"commercialIncomes,omitempty"`
	FacilityIncomes     []FacilityIncome     `bson:"facilityIncomes,omitempty" json:"facilityIncomes,omitempty"`
	Expenses            []Expense            `bson:"expenses,omitempty" json:"expenses,omitempty"`
	PurchaseFees        []PurchaseFee        `bson:"purchaseFees,omitempty" json:"purchaseFees,omitempty"`
	CapitalImprovements []CapitalImprovement `bson:"capitalImprovements,omitempty" json:"capitalImprovements,omitempty"`
	AnnualProjections   []AnnualProjection   `bson:"annualProjections,omitempty" json:"annualProjections,omitempty"`
}

// RentalIncome represents income from rental units
type RentalIncome struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText             string             `bson:"nameText" json:"nameText"`
	MonthlyAmount        decimal.Decimal    `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount         decimal.Decimal    `bson:"annualAmount" json:"annualAmount"`
	MonthlyAmountPerUnit decimal.Decimal    `bson:"monthlyAmountPerUnit" json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  decimal.Decimal    `bson:"annualAmountPerUnit" json:"annualAmountPerUnit"`
	NumberOfUnits        decimal.Decimal    `bson:"numberOfUnits" json:"numberOfUnits"`
	Frequency            decimal.Decimal    `bson:"frequency" json:"frequency"`
	TypeID               int                `bson:"typeId" json:"typeId"`
}

// CommercialIncome represents income from commercial spaces
type CommercialIncome struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText             string             `bson:"nameText" json:"nameText"`
	MonthlyAmount        decimal.Decimal    `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount         decimal.Decimal    `bson:"annualAmount" json:"annualAmount"`
	MonthlyAmountPerUnit decimal.Decimal    `bson:"monthlyAmountPerUnit" json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  decimal.Decimal    `bson:"annualAmountPerUnit" json:"annualAmountPerUnit"`
	AreaInSquareFeet     decimal.Decimal    `bson:"areaInSquareFeet" json:"areaInSquareFeet"`
	UnitType             string             `bson:"unitType" json:"unitType"`
	UnitValue            decimal.Decimal    `bson:"unitValue" json:"unitValue"`
	Frequency            decimal.Decimal    `bson:"frequency" json:"frequency"`
	TypeID               int                `bson:"typeId" json:"typeId"`
}

// FacilityIncome represents income from property facilities
type FacilityIncome struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText      string             `bson:"nameText" json:"nameText"`
	MonthlyAmount decimal.Decimal    `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount  decimal.Decimal    `bson:"annualAmount" json:"annualAmount"`
	Frequency     decimal.Decimal    `bson:"frequency" json:"frequency"`
	TypeID        int                `bson:"typeId" json:"typeId"`
}

// Expense represents property expenses
type Expense struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText      string             `bson:"nameText" json:"nameText"`
	MonthlyAmount decimal.Decimal    `bson:"monthlyAmount" json:"monthlyAmount"`
	AnnualAmount  decimal.Decimal    `bson:"annualAmount" json:"annualAmount"`
	Frequency     decimal.Decimal    `bson:"frequency" json:"frequency"`
	Percent       decimal.Decimal    `bson:"percent" json:"percent"`
	TypeID        int                `bson:"typeId" json:"typeId"`
}

// ExpenseCategory represents categories for expenses
type ExpenseCategory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryName string             `bson:"categoryName" json:"categoryName"`
}

// FacilityIncomeCategory represents categories for facility income
type FacilityIncomeCategory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryName string             `bson:"categoryName" json:"categoryName"`
}

// PurchaseFee represents fees associated with property purchase
type PurchaseFee struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText string             `bson:"nameText" json:"nameText"`
	Amount   decimal.Decimal    `bson:"amount" json:"amount"`
	TypeID   int                `bson:"typeId" json:"typeId"`
}

// CapitalImprovement represents capital improvements to property
type CapitalImprovement struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NameText string             `bson:"nameText" json:"nameText"`
	Amount   decimal.Decimal    `bson:"amount" json:"amount"`
	TypeID   int                `bson:"typeId" json:"typeId"`
}

// AnnualProjection represents annual financial projections
type AnnualProjection struct {
	ID                                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Year                                decimal.Decimal    `bson:"year" json:"year"`
	CashFlow                            decimal.Decimal    `bson:"cashFlow" json:"cashFlow"`
	DebtRemaining                       decimal.Decimal    `bson:"debtRemaining" json:"debtRemaining"`
	SalesPrice                          decimal.Decimal    `bson:"salesPrice" json:"salesPrice"`
	LegalFees                           decimal.Decimal    `bson:"legalFees" json:"legalFees"`
	ProceedsOfSale                      decimal.Decimal    `bson:"proceedsOfSale" json:"proceedsOfSale"`
	TotalReturn                         decimal.Decimal    `bson:"totalReturn" json:"totalReturn"`
	InitialInvestment                   decimal.Decimal    `bson:"initialInvestment" json:"initialInvestment"`
	ReturnOnInvestmentRate              decimal.Decimal    `bson:"returnOnInvestmentRate" json:"returnOnInvestmentRate"`
	ReturnOnInvestmentPercent           decimal.Decimal    `bson:"returnOnInvestmentPercent" json:"returnOnInvestmentPercent"`
	AnnualizedReturnOnInvestmentRate    decimal.Decimal    `bson:"annualizedReturnOnInvestmentRate" json:"annualizedReturnOnInvestmentRate"`
	AnnualizedReturnOnInvestmentPercent decimal.Decimal    `bson:"annualizedReturnOnInvestmentPercent" json:"annualizedReturnOnInvestmentPercent"`
}
