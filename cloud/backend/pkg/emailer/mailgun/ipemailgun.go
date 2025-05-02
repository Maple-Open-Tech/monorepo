package mailgun

import (
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
)

// NewIncomePropertyEvaluatorModuleEmailer creates a new emailer for the Income Property Evaluator module.
func NewIncomePropertyEvaluatorModuleEmailer(cfg *config.Configuration) Emailer {
	emailerConfigProvider := NewMailgunConfigurationProvider(
		cfg.IPEMailgun.SenderEmail,
		cfg.IPEMailgun.Domain,
		cfg.IPEMailgun.APIBase,
		cfg.IPEMailgun.MaintenanceEmail,
		cfg.IPEMailgun.FrontendDomain,
		cfg.IPEMailgun.BackendDomain,
		cfg.IPEMailgun.APIKey,
	)

	return NewEmailer(emailerConfigProvider)
}
