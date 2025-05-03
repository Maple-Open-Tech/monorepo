package mortgage

import (
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mortgage represents mortgage details for a property
type Mortgage struct {
	ID                                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FinancialAnalysisID                primitive.ObjectID `bson:"financialAnalysisId" json:"financialAnalysisId"`
	LoanAmount                         decimal.Decimal    `bson:"loanAmount" json:"loanAmount"`
	LoanPurchaseAmount                 decimal.Decimal    `bson:"loanPurchaseAmount" json:"loanPurchaseAmount"`
	DownPayment                        decimal.Decimal    `bson:"downPayment" json:"downPayment"`
	AnnualInterestRate                 decimal.Decimal    `bson:"annualInterestRate" json:"annualInterestRate"`
	AmortizationYear                   decimal.Decimal    `bson:"amortizationYear" json:"amortizationYear"`
	PaymentFrequency                   decimal.Decimal    `bson:"paymentFrequency" json:"paymentFrequency"`
	CompoundingPeriod                  decimal.Decimal    `bson:"compoundingPeriod" json:"compoundingPeriod"`
	FirstPaymentDate                   time.Time          `bson:"firstPaymentDate" json:"firstPaymentDate"`
	Insurance                          string             `bson:"insurance" json:"insurance"`
	InsuranceAmount                    decimal.Decimal    `bson:"insuranceAmount" json:"insuranceAmount"`
	MortgagePaymentPerPaymentFrequency decimal.Decimal    `bson:"mortgagePaymentPerPaymentFrequency" json:"mortgagePaymentPerPaymentFrequency"`
	InterestRatePerPaymentFrequency    decimal.Decimal    `bson:"interestRatePerPaymentFrequency" json:"interestRatePerPaymentFrequency"`
	TotalNumberOfPaymentsPerFrequency  decimal.Decimal    `bson:"totalNumberOfPaymentsPerFrequency" json:"totalNumberOfPaymentsPerFrequency"`
	PercentFinanced                    decimal.Decimal    `bson:"percentFinanced" json:"percentFinanced"`

	// Embedded payment schedule
	MortgagePaymentSchedule []MortgageInterval `bson:"mortgagePaymentSchedule,omitempty" json:"mortgagePaymentSchedule,omitempty"`
}

// MortgageInterval represents a payment in the mortgage schedule
type MortgageInterval struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Interval            decimal.Decimal    `bson:"interval" json:"interval"`
	PaymentDate         time.Time          `bson:"paymentDate" json:"paymentDate"`
	PaymentAmount       decimal.Decimal    `bson:"paymentAmount" json:"paymentAmount"`
	PrincipleAmount     decimal.Decimal    `bson:"principleAmount" json:"principleAmount"`
	InterestAmount      decimal.Decimal    `bson:"interestAmount" json:"interestAmount"`
	LoanBalance         decimal.Decimal    `bson:"loanBalance" json:"loanBalance"`
	Year                decimal.Decimal    `bson:"year" json:"year"`
	TotalPaidToBank     decimal.Decimal    `bson:"totalPaidToBank" json:"totalPaidToBank"`
	TotalPaidToInterest decimal.Decimal    `bson:"totalPaidToInterest" json:"totalPaidToInterest"`
}
