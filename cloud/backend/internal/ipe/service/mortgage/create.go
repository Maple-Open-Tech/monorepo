// cloud/backend/internal/ipe/service/mortgage/create.go
package mortgage

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	uc_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type CreateMortgageRequestDTO struct {
	FinancialAnalysisID primitive.ObjectID `json:"financialAnalysisId"`
	LoanAmount          string             `json:"loanAmount"`
	DownPayment         string             `json:"downPayment"`
	AnnualInterestRate  string             `json:"annualInterestRate"`
	AmortizationYear    string             `json:"amortizationYear"`
	PaymentFrequency    string             `json:"paymentFrequency"`
	CompoundingPeriod   string             `json:"compoundingPeriod"`
	FirstPaymentDate    string             `json:"firstPaymentDate"`
	Insurance           string             `json:"insurance"`
	InsuranceAmount     string             `json:"insuranceAmount"`
}

type CreateMortgageResponseDTO struct {
	ID                  primitive.ObjectID `json:"id"`
	FinancialAnalysisID primitive.ObjectID `json:"financialAnalysisId"`
	LoanAmount          string             `json:"loanAmount"`
	DownPayment         string             `json:"downPayment"`
	AnnualInterestRate  string             `json:"annualInterestRate"`
	AmortizationYear    string             `json:"amortizationYear"`
	PaymentFrequency    string             `json:"paymentFrequency"`
	CompoundingPeriod   string             `json:"compoundingPeriod"`
	FirstPaymentDate    string             `json:"firstPaymentDate"`
	Insurance           string             `json:"insurance"`
	InsuranceAmount     string             `json:"insuranceAmount"`
}

type CreateMortgageService interface {
	Execute(ctx context.Context, request *CreateMortgageRequestDTO) (*CreateMortgageResponseDTO, error)
}

type createMortgageServiceImpl struct {
	config                *config.Configuration
	logger                *zap.Logger
	mortgageCreateUseCase uc_mortgage.MortgageCreateUseCase
}

func NewCreateMortgageService(
	config *config.Configuration,
	logger *zap.Logger,
	mortgageCreateUseCase uc_mortgage.MortgageCreateUseCase,
) CreateMortgageService {
	return &createMortgageServiceImpl{
		config:                config,
		logger:                logger,
		mortgageCreateUseCase: mortgageCreateUseCase,
	}
}

func (svc *createMortgageServiceImpl) Execute(ctx context.Context, req *CreateMortgageRequestDTO) (*CreateMortgageResponseDTO, error) {
	// Validate request
	if req == nil {
		return nil, httperror.NewForBadRequestWithSingleField("request", "Request is required")
	}

	errors := make(map[string]string)
	if req.FinancialAnalysisID.IsZero() {
		errors["financialAnalysisId"] = "Financial analysis ID is required"
	}
	if req.LoanAmount == "" {
		errors["loanAmount"] = "Loan amount is required"
	}
	if req.AnnualInterestRate == "" {
		errors["annualInterestRate"] = "Annual interest rate is required"
	}
	if req.AmortizationYear == "" {
		errors["amortizationYear"] = "Amortization year is required"
	}

	if len(errors) > 0 {
		return nil, httperror.NewForBadRequest(&errors)
	}

	// Parse decimal values
	loanAmount, err := decimal.NewFromString(req.LoanAmount)
	if err != nil {
		return nil, httperror.NewForBadRequestWithSingleField("loanAmount", "Invalid loan amount format")
	}

	downPayment := decimal.Zero
	if req.DownPayment != "" {
		downPayment, err = decimal.NewFromString(req.DownPayment)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("downPayment", "Invalid down payment format")
		}
	}

	annualInterestRate, err := decimal.NewFromString(req.AnnualInterestRate)
	if err != nil {
		return nil, httperror.NewForBadRequestWithSingleField("annualInterestRate", "Invalid annual interest rate format")
	}

	amortizationYear, err := decimal.NewFromString(req.AmortizationYear)
	if err != nil {
		return nil, httperror.NewForBadRequestWithSingleField("amortizationYear", "Invalid amortization year format")
	}

	paymentFrequency := decimal.NewFromInt(12) // Monthly by default
	if req.PaymentFrequency != "" {
		paymentFrequency, err = decimal.NewFromString(req.PaymentFrequency)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("paymentFrequency", "Invalid payment frequency format")
		}
	}

	compoundingPeriod := decimal.NewFromInt(2) // Semi-annually by default (Canadian standard)
	if req.CompoundingPeriod != "" {
		compoundingPeriod, err = decimal.NewFromString(req.CompoundingPeriod)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("compoundingPeriod", "Invalid compounding period format")
		}
	}

	firstPaymentDate := time.Now().AddDate(0, 1, 0) // One month from now by default
	if req.FirstPaymentDate != "" {
		firstPaymentDate, err = time.Parse("2006-01-02", req.FirstPaymentDate)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("firstPaymentDate", "Invalid date format, use YYYY-MM-DD")
		}
	}

	insuranceAmount := decimal.Zero
	if req.InsuranceAmount != "" {
		insuranceAmount, err = decimal.NewFromString(req.InsuranceAmount)
		if err != nil {
			return nil, httperror.NewForBadRequestWithSingleField("insuranceAmount", "Invalid insurance amount format")
		}
	}

	// Calculate loan purchase amount (loan + down payment)
	loanPurchaseAmount := loanAmount.Add(downPayment)

	// Calculate percent financed
	percentFinanced := decimal.Zero
	if !loanPurchaseAmount.IsZero() {
		percentFinanced = loanAmount.Div(loanPurchaseAmount).Mul(decimal.NewFromInt(100))
	}

	// Create mortgage domain object
	mortgage := &dom_mortgage.Mortgage{
		ID:                  primitive.NewObjectID(),
		FinancialAnalysisID: req.FinancialAnalysisID,
		LoanAmount:          loanAmount,
		LoanPurchaseAmount:  loanPurchaseAmount,
		DownPayment:         downPayment,
		AnnualInterestRate:  annualInterestRate,
		AmortizationYear:    amortizationYear,
		PaymentFrequency:    paymentFrequency,
		CompoundingPeriod:   compoundingPeriod,
		FirstPaymentDate:    firstPaymentDate,
		Insurance:           req.Insurance,
		InsuranceAmount:     insuranceAmount,
		PercentFinanced:     percentFinanced,
		// These calculations would typically be more complex in a real system
		MortgagePaymentPerPaymentFrequency: decimal.Zero,
		InterestRatePerPaymentFrequency:    annualInterestRate.Div(paymentFrequency),
		TotalNumberOfPaymentsPerFrequency:  amortizationYear.Mul(paymentFrequency),
		MortgagePaymentSchedule:            []dom_mortgage.MortgageInterval{},
	}

	// Create mortgage
	id, err := svc.mortgageCreateUseCase.Execute(ctx, mortgage)
	if err != nil {
		svc.logger.Error("Failed to create mortgage", zap.Error(err))
		return nil, err
	}

	// Create response
	response := &CreateMortgageResponseDTO{
		ID:                  id,
		FinancialAnalysisID: mortgage.FinancialAnalysisID,
		LoanAmount:          mortgage.LoanAmount.String(),
		DownPayment:         mortgage.DownPayment.String(),
		AnnualInterestRate:  mortgage.AnnualInterestRate.String(),
		AmortizationYear:    mortgage.AmortizationYear.String(),
		PaymentFrequency:    mortgage.PaymentFrequency.String(),
		CompoundingPeriod:   mortgage.CompoundingPeriod.String(),
		FirstPaymentDate:    mortgage.FirstPaymentDate.Format("2006-01-02"),
		Insurance:           mortgage.Insurance,
		InsuranceAmount:     mortgage.InsuranceAmount.String(),
	}

	return response, nil
}
