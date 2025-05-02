package templatedemailer

import (
	"context"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/emailer/mailgun"
)

// TemplatedEmailer Is adapter for responsive HTML email templates sender.
type TemplatedEmailer interface {
	SendUserVerificationEmail(ctx context.Context, monolithModule int, email, verificationCode, firstName string) error
	SendUserPasswordResetEmail(ctx context.Context, monolithModule int, email, verificationCode, firstName string) error
}

type templatedEmailer struct {
	incomePropertyEmailer mailgun.Emailer
	maplesendEmailer      mailgun.Emailer
}

func NewTemplatedEmailer(incomePropertyEmailer mailgun.Emailer, maplesendEmailer mailgun.Emailer) TemplatedEmailer {
	return &templatedEmailer{
		incomePropertyEmailer: incomePropertyEmailer,
		maplesendEmailer:      maplesendEmailer,
	}
}
