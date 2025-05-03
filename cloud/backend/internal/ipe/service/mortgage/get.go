// cloud/backend/internal/ipe/service/mortgage/get.go
package mortgage

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
	uc_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/mortgage"
)

type MortgageResponseDTO struct {
	ID                                 primitive.ObjectID    `json:"id"`
	FinancialAnalysisID                primitive.ObjectID    `json:"financialAnalysisId"`
	LoanAmount                         string                `json:"loanAmount"`
	LoanPurchaseAmount                 string                `json:"loanPurchaseAmount"`
	DownPayment                        string                `json:"downPayment"`
	AnnualInterestRate                 string                `json:"annualInterestRate"`
	AmortizationYear                   string                `json:"amortizationYear"`
	PaymentFrequency                   string                `json:"paymentFrequency"`
	CompoundingPeriod                  string                `json:"compoundingPeriod"`
	FirstPaymentDate                   string                `json:"firstPaymentDate"`
	Insurance                          string                `json:"insurance"`
	InsuranceAmount                    string                `json:"insuranceAmount"`
	MortgagePaymentPerPaymentFrequency string                `json:"mortgagePaymentPerPaymentFrequency"`
	InterestRatePerPaymentFrequency    string                `json:"interestRatePerPaymentFrequency"`
	TotalNumberOfPaymentsPerFrequency  string                `json:"totalNumberOfPaymentsPerFrequency"`
	PercentFinanced                    string                `json:"percentFinanced"`
	MortgagePaymentSchedule            []MortgageIntervalDTO `json:"mortgagePaymentSchedule,omitempty"`
}

type MortgageIntervalDTO struct {
	ID                  primitive.ObjectID `json:"id"`
	Interval            string             `json:"interval"`
	PaymentDate         string             `json:"paymentDate"`
	PaymentAmount       string             `json:"paymentAmount"`
	PrincipleAmount     string             `json:"principleAmount"`
	InterestAmount      string             `json:"interestAmount"`
	LoanBalance         string             `json:"loanBalance"`
	Year                string             `json:"year"`
	TotalPaidToBank     string             `json:"totalPaidToBank"`
	TotalPaidToInterest string             `json:"totalPaidToInterest"`
}

type GetMortgageService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*MortgageResponseDTO, error)
	ExecuteByFinancialAnalysisID(ctx context.Context, financialAnalysisID primitive.ObjectID) (*MortgageResponseDTO, error)
}

type getMortgageServiceImpl struct {
	config                                  *config.Configuration
	logger                                  *zap.Logger
	mortgageGetByIDUseCase                  uc_mortgage.MortgageGetByIDUseCase
	mortgageGetByFinancialAnalysisIDUseCase uc_mortgage.MortgageGetByFinancialAnalysisIDUseCase
}

func NewGetMortgageService(
	config *config.Configuration,
	logger *zap.Logger,
	mortgageGetByIDUseCase uc_mortgage.MortgageGetByIDUseCase,
	mortgageGetByFinancialAnalysisIDUseCase uc_mortgage.MortgageGetByFinancialAnalysisIDUseCase,
) GetMortgageService {
	return &getMortgageServiceImpl{
		config:                                  config,
		logger:                                  logger,
		mortgageGetByIDUseCase:                  mortgageGetByIDUseCase,
		mortgageGetByFinancialAnalysisIDUseCase: mortgageGetByFinancialAnalysisIDUseCase,
	}
}

func (svc *getMortgageServiceImpl) Execute(ctx context.Context, id primitive.ObjectID) (*MortgageResponseDTO, error) {
	// Get mortgage
	mortgage, err := svc.mortgageGetByIDUseCase.Execute(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Mortgage not found", zap.String("id", id.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get mortgage", zap.Error(err))
		return nil, err
	}
	if mortgage == nil {
		svc.logger.Warn("Mortgage not found", zap.String("id", id.Hex()))
		return nil, errors.New("mortgage not found")
	}

	return svc.mapToResponseDTO(mortgage), nil
}

func (svc *getMortgageServiceImpl) ExecuteByFinancialAnalysisID(ctx context.Context, financialAnalysisID primitive.ObjectID) (*MortgageResponseDTO, error) {
	// Get mortgage by financial analysis ID
	mortgage, err := svc.mortgageGetByFinancialAnalysisIDUseCase.Execute(ctx, financialAnalysisID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Warn("Mortgage not found for financial analysis", zap.String("financialAnalysisId", financialAnalysisID.Hex()))
			return nil, err
		}
		svc.logger.Error("Failed to get mortgage by financial analysis ID", zap.Error(err))
		return nil, err
	}
	if mortgage == nil {
		svc.logger.Warn("Mortgage not found for financial analysis", zap.String("financialAnalysisId", financialAnalysisID.Hex()))
		return nil, errors.New("mortgage not found")
	}

	return svc.mapToResponseDTO(mortgage), nil
}

func (svc *getMortgageServiceImpl) mapToResponseDTO(mortgage *dom_mortgage.Mortgage) *MortgageResponseDTO {
	// Map payment schedule
	paymentSchedule := make([]MortgageIntervalDTO, len(mortgage.MortgagePaymentSchedule))
	for i, interval := range mortgage.MortgagePaymentSchedule {
		paymentSchedule[i] = MortgageIntervalDTO{
			ID:                  interval.ID,
			Interval:            interval.Interval.String(),
			PaymentDate:         interval.PaymentDate.Format("2006-01-02"),
			PaymentAmount:       interval.PaymentAmount.String(),
			PrincipleAmount:     interval.PrincipleAmount.String(),
			InterestAmount:      interval.InterestAmount.String(),
			LoanBalance:         interval.LoanBalance.String(),
			Year:                interval.Year.String(),
			TotalPaidToBank:     interval.TotalPaidToBank.String(),
			TotalPaidToInterest: interval.TotalPaidToInterest.String(),
		}
	}

	// Create response
	response := &MortgageResponseDTO{
		ID:                                 mortgage.ID,
		FinancialAnalysisID:                mortgage.FinancialAnalysisID,
		LoanAmount:                         mortgage.LoanAmount.String(),
		LoanPurchaseAmount:                 mortgage.LoanPurchaseAmount.String(),
		DownPayment:                        mortgage.DownPayment.String(),
		AnnualInterestRate:                 mortgage.AnnualInterestRate.String(),
		AmortizationYear:                   mortgage.AmortizationYear.String(),
		PaymentFrequency:                   mortgage.PaymentFrequency.String(),
		CompoundingPeriod:                  mortgage.CompoundingPeriod.String(),
		FirstPaymentDate:                   mortgage.FirstPaymentDate.Format("2006-01-02"),
		Insurance:                          mortgage.Insurance,
		InsuranceAmount:                    mortgage.InsuranceAmount.String(),
		MortgagePaymentPerPaymentFrequency: mortgage.MortgagePaymentPerPaymentFrequency.String(),
		InterestRatePerPaymentFrequency:    mortgage.InterestRatePerPaymentFrequency.String(),
		TotalNumberOfPaymentsPerFrequency:  mortgage.TotalNumberOfPaymentsPerFrequency.String(),
		PercentFinanced:                    mortgage.PercentFinanced.String(),
		MortgagePaymentSchedule:            paymentSchedule,
	}

	return response
}
