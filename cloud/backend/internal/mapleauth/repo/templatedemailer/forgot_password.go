package templatedemailer

import (
	"context"
)

func (impl *templatedEmailer) SendUserPasswordResetEmail(ctx context.Context, monolithModule int, email, verificationCode, firstName string) error {

	return nil
}
