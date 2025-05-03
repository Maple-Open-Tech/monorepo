// cloud/backend/internal/ipe/service/financialanalysis/get.go
package financialanalysis

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	uc_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/financialanalysis"
	uc_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/mortgage"
)

type FinancialAnalysisResponseDTO struct {
	ID                     primitive.ObjectID    `json:"id"`
	PropertyID             primitive.ObjectID    `json:"propertyId"`
	PurchasePrice          string                `json:"purchasePrice"`
	AnnualGrossIncome      string                `json:"annualGrossIncome"`
	MonthlyGrossIncome     string                `json:"monthlyGrossIncome"`
	AnnualExpense          string                `json:"annualExpense"`
	MonthlyExpense         string                `json:"monthlyExpense"`
	AnnualNetIncome        string                `json:"annualNetIncome"`
	MonthlyNetIncome       string                `json:"monthlyNetIncome"`
	AnnualCashFlow         string                `json:"annualCashFlow"`
	MonthlyCashFlow        string                `json:"monthlyCashFlow"`
	CapRateWithMortgage    string                `json:"capRateWithMortgage"`
	CapRateWithoutMortgage string                `json:"capRateWithoutMortgage"`
	RentalIncomes          []RentalIncomeDTO     `json:"rentalIncomes,omitempty"`
	CommercialIncomes      []CommercialIncomeDTO `json:"commercialIncomes,omitempty"`
	FacilityIncomes        []FacilityIncomeDTO   `json:"facilityIncomes,omitempty"`
	Expenses               []ExpenseDTO          `json:"expenses,omitempty"`
	HasMortgage            bool                  `json:"hasMortgage"`
}

type RentalIncomeDTO struct {
	ID                   primitive.ObjectID `json:"id"`
	NameText             string             `json:"nameText"`
	MonthlyAmount        string             `json:"monthlyAmount"`
	AnnualAmount         string             `json:"annualAmount"`
	MonthlyAmountPerUnit string             `json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  string             `json:"annualAmountPerUnit"`
	NumberOfUnits        string             `json:"numberOfUnits"`
	Frequency            string             `json:"frequency"`
	TypeID               int                `json:"typeId"`
}

type CommercialIncomeDTO struct {
	ID                   primitive.ObjectID `json:"id"`
	NameText             string             `json:"nameText"`
	MonthlyAmount        string             `json:"monthlyAmount"`
	AnnualAmount         string             `json:"annualAmount"`
	MonthlyAmountPerUnit string             `json:"monthlyAmountPerUnit"`
	AnnualAmountPerUnit  string             `json:"annualAmountPerUnit"`
	AreaInSquareFeet     string             `json:"areaInSquareFeet"`
	UnitType             string             `json:"unitType"`
	UnitValue            string             `json:"unitValue"`
	Frequency            string             `json:"frequency"`
	TypeID               int                `json:"typeId"`
}

type FacilityIncomeDTO struct {
	ID            primitive.ObjectID `json:"id"`
	NameText      string             `json:"nameText"`
	MonthlyAmount string             `json:"monthlyAmount"`
	AnnualAmount  string             `json:"annualAmount"`
	Frequency     string             `json:"frequency"`
	TypeID        int                `json:"typeId"`
}

type ExpenseDTO struct {
	ID            primitive.ObjectID `json:"id"`
	NameText      string             `json:"nameText"`
	MonthlyAmount string             `json:"monthlyAmount"`
	AnnualAmount  string             `json:"annualAmount"`
	Frequency     string             `json:"frequency"`
	Percent       string             `json:"percent"`
	TypeID        int                `json:"typeId"`
}

type GetFinancialAnalysisService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*FinancialAnalysisResponseDTO, error)
	ExecuteByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*FinancialAnalysisResponseDTO, error)
}

type getFinancialAnalysisServiceImpl struct {
	config                                  *config.Configuration
	logger                                  *zap.Logger
	financialAnalysisGetByIDUseCase         uc_financial.FinancialAnalysisGetByIDUseCase
	financialAnalysisGetByPropertyIDUseCase uc_financial.FinancialAnalysisGetByPropertyIDUseCase
	mortgageGetByFinancialAnalysisIDUseCase uc_mortgage.MortgageGetByFinancialAnalysisIDUseCase
}

func NewGetFinancialAnalysisService(
	config *config.Configuration,
	logger *zap.Logger,
	financialAnalysisGetByIDUseCase uc_financial.FinancialAnalysisGetByIDUseCase,
	financialAnalysisGetByPropertyIDUseCase uc_financial.FinancialAnalysisGetByPropertyIDUseCase,
	mortgageGetByFinancialAnalysisIDUseCase uc_mortgage.MortgageGetByFinancialAnalysisIDUseCase,
) GetFinancialAnalysisService {
	return &getFinancialAnalysisServiceImpl{
		config:                                  config,
		logger:                                  logger,
		financialAnalysisGetByIDUseCase:         financialAnalysisGetByIDUseCase,
		financialAnalysisGetByPropertyIDUseCase: financialAnalysisGetByPropertyIDUseCase,
		mortgageGetByFinancialAnalysisIDUseCase: mortgageGetByFinancialAnalysisIDUseCase,
	}
}

func (svc *getFinancialAnalysisServiceImpl) Execute(ctx context.Context, id primitive.ObjectID) (*FinancialAnalysisResponseDTO, error) {
	// Get financial analysis
	analysis, err := svc.financialAnalysisGetByIDUseCase.Execute(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Financial analysis not found", zap.String("id", id.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get financial analysis", zap.Error(err))
		return nil, err
	}
	if analysis == nil {
		svc.logger.Warn("Financial analysis not found", zap.String("id", id.Hex()))
		return nil, errors.New("financial analysis not found")
	}

	return svc.mapToResponseDTO(ctx, analysis)
}

func (svc *getFinancialAnalysisServiceImpl) ExecuteByPropertyID(ctx context.Context, propertyID primitive.ObjectID) (*FinancialAnalysisResponseDTO, error) {
	// Get financial analysis by property ID
	analysis, err := svc.financialAnalysisGetByPropertyIDUseCase.Execute(ctx, propertyID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Financial analysis not found for property", zap.String("propertyId", propertyID.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get financial analysis by property ID", zap.Error(err))
		return nil, err
	}
	if analysis == nil {
		svc.logger.Warn("Financial analysis not found for property", zap.String("propertyId", propertyID.Hex()))
		return nil, errors.New("financial analysis not found")
	}

	return svc.mapToResponseDTO(ctx, analysis)
}

func (svc *getFinancialAnalysisServiceImpl) mapToResponseDTO(ctx context.Context, analysis *dom_financial.FinancialAnalysis) (*FinancialAnalysisResponseDTO, error) {
	// Check if mortgage exists
	mortgage, err := svc.mortgageGetByFinancialAnalysisIDUseCase.Execute(ctx, analysis.ID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		svc.logger.Error("Failed to check for mortgage", zap.Error(err))
		return nil, err
	}
	hasMortgage := mortgage != nil

	// Map rental incomes
	rentalIncomes := make([]RentalIncomeDTO, len(analysis.RentalIncomes))
	for i, income := range analysis.RentalIncomes {
		rentalIncomes[i] = RentalIncomeDTO{
			ID:                   income.ID,
			NameText:             income.NameText,
			MonthlyAmount:        income.MonthlyAmount.String(),
			AnnualAmount:         income.AnnualAmount.String(),
			MonthlyAmountPerUnit: income.MonthlyAmountPerUnit.String(),
			AnnualAmountPerUnit:  income.AnnualAmountPerUnit.String(),
			NumberOfUnits:        income.NumberOfUnits.String(),
			Frequency:            income.Frequency.String(),
			TypeID:               income.TypeID,
		}
	}

	// Map commercial incomes
	commercialIncomes := make([]CommercialIncomeDTO, len(analysis.CommercialIncomes))
	for i, income := range analysis.CommercialIncomes {
		commercialIncomes[i] = CommercialIncomeDTO{
			ID:                   income.ID,
			NameText:             income.NameText,
			MonthlyAmount:        income.MonthlyAmount.String(),
			AnnualAmount:         income.AnnualAmount.String(),
			MonthlyAmountPerUnit: income.MonthlyAmountPerUnit.String(),
			AnnualAmountPerUnit:  income.AnnualAmountPerUnit.String(),
			AreaInSquareFeet:     income.AreaInSquareFeet.String(),
			UnitType:             income.UnitType,
			UnitValue:            income.UnitValue.String(),
			Frequency:            income.Frequency.String(),
			TypeID:               income.TypeID,
		}
	}

	// Map facility incomes
	facilityIncomes := make([]FacilityIncomeDTO, len(analysis.FacilityIncomes))
	for i, income := range analysis.FacilityIncomes {
		facilityIncomes[i] = FacilityIncomeDTO{
			ID:            income.ID,
			NameText:      income.NameText,
			MonthlyAmount: income.MonthlyAmount.String(),
			AnnualAmount:  income.AnnualAmount.String(),
			Frequency:     income.Frequency.String(),
			TypeID:        income.TypeID,
		}
	}

	// Map expenses
	expenses := make([]ExpenseDTO, len(analysis.Expenses))
	for i, expense := range analysis.Expenses {
		expenses[i] = ExpenseDTO{
			ID:            expense.ID,
			NameText:      expense.NameText,
			MonthlyAmount: expense.MonthlyAmount.String(),
			AnnualAmount:  expense.AnnualAmount.String(),
			Frequency:     expense.Frequency.String(),
			Percent:       expense.Percent.String(),
			TypeID:        expense.TypeID,
		}
	}

	// Create response
	response := &FinancialAnalysisResponseDTO{
		ID:                     analysis.ID,
		PropertyID:             analysis.PropertyID,
		PurchasePrice:          analysis.PurchasePrice.String(),
		AnnualGrossIncome:      analysis.AnnualGrossIncome.String(),
		MonthlyGrossIncome:     analysis.MonthlyGrossIncome.String(),
		AnnualExpense:          analysis.AnnualExpense.String(),
		MonthlyExpense:         analysis.MonthlyExpense.String(),
		AnnualNetIncome:        analysis.AnnualNetIncome.String(),
		MonthlyNetIncome:       analysis.MonthlyNetIncome.String(),
		AnnualCashFlow:         analysis.AnnualCashFlow.String(),
		MonthlyCashFlow:        analysis.MonthlyCashFlow.String(),
		CapRateWithMortgage:    analysis.CapRateWithMortgage.String(),
		CapRateWithoutMortgage: analysis.CapRateWithoutMortgage.String(),
		RentalIncomes:          rentalIncomes,
		CommercialIncomes:      commercialIncomes,
		FacilityIncomes:        facilityIncomes,
		Expenses:               expenses,
		HasMortgage:            hasMortgage,
	}

	return response, nil
}
