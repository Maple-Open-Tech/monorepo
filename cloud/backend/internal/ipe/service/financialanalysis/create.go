// cloud/backend/internal/ipe/service/financialanalysis/create.go
package financialanalysis

import (
	"context"

	"go.uber.org/zap"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
	uc_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateFinancialAnalysisRequestDTO struct {
	PropertyID     primitive.ObjectID `json:"propertyId"`
	PurchasePrice  string             `json:"purchasePrice"`
	InflationRate  string             `json:"inflationRate"`
	BuyingFeeRate  string             `json:"buyingFeeRate"`
	SellingFeeRate string             `json:"sellingFeeRate"`
}

type CreateFinancialAnalysisResponseDTO struct {
	ID             primitive.ObjectID `json:"id"`
	PropertyID     primitive.ObjectID `json:"propertyId"`
	PurchasePrice  string             `json:"purchasePrice"`
	InflationRate  string             `json:"inflationRate"`
	BuyingFeeRate  string             `json:"buyingFeeRate"`
	SellingFeeRate string             `json:"sellingFeeRate"`
}

type CreateFinancialAnalysisService interface {
	Execute(ctx context.Context, request *CreateFinancialAnalysisRequestDTO) (*CreateFinancialAnalysisResponseDTO, error)
}

type createFinancialAnalysisServiceImpl struct {
	config                         *config.Configuration
	logger                         *zap.Logger
	financialAnalysisCreateUseCase uc_financial.FinancialAnalysisCreateUseCase
}

func NewCreateFinancialAnalysisService(
	config *config.Configuration,
	logger *zap.Logger,
	financialAnalysisCreateUseCase uc_financial.FinancialAnalysisCreateUseCase,
) CreateFinancialAnalysisService {
	return &createFinancialAnalysisServiceImpl{
		config:                         config,
		logger:                         logger,
		financialAnalysisCreateUseCase: financialAnalysisCreateUseCase,
	}
}

func (svc *createFinancialAnalysisServiceImpl) Execute(ctx context.Context, req *CreateFinancialAnalysisRequestDTO) (*CreateFinancialAnalysisResponseDTO, error) {
	// Validate request
	if req == nil {
		return nil, httperror.NewForBadRequestWithSingleField("request", "Request is required")
	}

	errors := make(map[string]string)
	if req.PropertyID.IsZero() {
		errors["propertyId"] = "Property ID is required"
	}
	if req.PurchasePrice == "" {
		errors["purchasePrice"] = "Purchase price is required"
	}

	if len(errors) > 0 {
		return nil, httperror.NewForBadRequest(&errors)
	}

	// Parse decimal values
	purchasePrice, err := decimal.NewFromString(req.PurchasePrice)
	if err != nil {
		return nil, httperror.NewForBadRequestWithSingleField("purchasePrice", "Invalid purchase price format")
	}

	inflationRate := decimal.NewFromFloat(2.0) // Default 2%
	if req.InflationRate != "" {
		inflationRate, err = decimal.NewFromString(req.InflationRate)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("inflationRate", "Invalid inflation rate format")
		}
	}

	buyingFeeRate := decimal.NewFromFloat(1.5) // Default 1.5%
	if req.BuyingFeeRate != "" {
		buyingFeeRate, err = decimal.NewFromString(req.BuyingFeeRate)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("buyingFeeRate", "Invalid buying fee rate format")
		}
	}

	sellingFeeRate := decimal.NewFromFloat(5.0) // Default 5%
	if req.SellingFeeRate != "" {
		sellingFeeRate, err = decimal.NewFromString(req.SellingFeeRate)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("sellingFeeRate", "Invalid selling fee rate format")
		}
	}

	// Create financial analysis domain object
	analysis := &dom_financial.FinancialAnalysis{
		ID:                        primitive.NewObjectID(),
		PropertyID:                req.PropertyID,
		PurchasePrice:             purchasePrice,
		InflationRate:             inflationRate,
		BuyingFeeRate:             buyingFeeRate,
		SellingFeeRate:            sellingFeeRate,
		AnnualGrossIncome:         decimal.Zero,
		MonthlyGrossIncome:        decimal.Zero,
		AnnualExpense:             decimal.Zero,
		MonthlyExpense:            decimal.Zero,
		AnnualNetIncome:           decimal.Zero,
		MonthlyNetIncome:          decimal.Zero,
		AnnualCashFlow:            decimal.Zero,
		MonthlyCashFlow:           decimal.Zero,
		CapRateWithMortgage:       decimal.Zero,
		CapRateWithoutMortgage:    decimal.Zero,
		AnnualRentalIncome:        decimal.Zero,
		MonthlyRentalIncome:       decimal.Zero,
		AnnualFacilityIncome:      decimal.Zero,
		MonthlyFacilityIncome:     decimal.Zero,
		CapitalImprovementsAmount: decimal.Zero,
		PurchaseFeesAmount:        decimal.Zero,
		InitialInvestmentAmount:   decimal.Zero,
		RentalIncomes:             []dom_financial.RentalIncome{},
		CommercialIncomes:         []dom_financial.CommercialIncome{},
		FacilityIncomes:           []dom_financial.FacilityIncome{},
		Expenses:                  []dom_financial.Expense{},
		PurchaseFees:              []dom_financial.PurchaseFee{},
		CapitalImprovements:       []dom_financial.CapitalImprovement{},
		AnnualProjections:         []dom_financial.AnnualProjection{},
	}

	// Create financial analysis
	id, err := svc.financialAnalysisCreateUseCase.Execute(ctx, analysis)
	if err != nil {
		svc.logger.Error("Failed to create financial analysis", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateFinancialAnalysisResponseDTO{
		ID:             id,
		PropertyID:     analysis.PropertyID,
		PurchasePrice:  analysis.PurchasePrice.String(),
		InflationRate:  analysis.InflationRate.String(),
		BuyingFeeRate:  analysis.BuyingFeeRate.String(),
		SellingFeeRate: analysis.SellingFeeRate.String(),
	}

	return response, nil
}
