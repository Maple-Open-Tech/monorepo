// cloud/backend/internal/ipe/usecase/module.go
package usecase

import (
	"go.uber.org/fx"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/bannedipaddress"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/emailer"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/evaluation"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/financialanalysis"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/incomeproperty"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/mortgage"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/client"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/owner"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/person/presenter"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/usecase/user"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			// Existing use cases
			bannedipaddress.NewCreateBannedIPAddressUseCase,
			bannedipaddress.NewBannedIPAddressListAllValuesUseCase,
			emailer.NewSendUserPasswordResetEmailUseCase,
			emailer.NewSendUserVerificationEmailUseCase,
			user.NewUserGetBySessionIDUseCase,
			user.NewUserCountByFilterUseCase,
			user.NewUserCreateUseCase,
			user.NewUserDeleteUserByEmailUseCase,
			user.NewUserDeleteByIDUseCase,
			user.NewUserGetByEmailUseCase,
			user.NewUserGetByIDUseCase,
			user.NewUserGetByVerificationCodeUseCase,
			user.NewUserListAllUseCase,
			user.NewUserListByFilterUseCase,
			user.NewUserUpdateUseCase,

			// Income Property use cases
			incomeproperty.NewIncomePropertyCreateUseCase,
			incomeproperty.NewIncomePropertyGetByIDUseCase,
			incomeproperty.NewIncomePropertyUpdateUseCase,
			incomeproperty.NewIncomePropertyDeleteUseCase,
			incomeproperty.NewIncomePropertyListAllUseCase,
			incomeproperty.NewIncomePropertyFindByCityUseCase,

			// Financial Analysis use cases
			financialanalysis.NewFinancialAnalysisCreateUseCase,
			financialanalysis.NewFinancialAnalysisGetByIDUseCase,
			financialanalysis.NewFinancialAnalysisGetByPropertyIDUseCase,
			financialanalysis.NewFinancialAnalysisUpdateUseCase,
			financialanalysis.NewAddRentalIncomeUseCase,

			// Mortgage use cases
			mortgage.NewMortgageCreateUseCase,
			mortgage.NewMortgageGetByIDUseCase,
			mortgage.NewMortgageGetByFinancialAnalysisIDUseCase,
			mortgage.NewAddMortgageIntervalUseCase,

			// Person use cases
			client.NewClientCreateUseCase,
			client.NewClientGetByIDUseCase,
			client.NewClientListAllUseCase,
			presenter.NewPresenterCreateUseCase,
			owner.NewOwnerCreateUseCase,

			// Evaluation use cases
			evaluation.NewEvaluationCreateUseCase,
			evaluation.NewEvaluationGetByIDUseCase,
			evaluation.NewEvaluationGetByPropertyIDUseCase,
			evaluation.NewEvaluationUpdateUseCase,
			evaluation.NewAddPropertyPhotoUseCase,
		),
	)
}
