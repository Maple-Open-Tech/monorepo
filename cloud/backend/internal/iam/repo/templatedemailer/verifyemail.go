package templatedemailer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path"
	"text/template"
)

func (impl *templatedEmailer) SendUserVerificationEmail(ctx context.Context, monolithModule int, email, verificationCode, firstName string) error {
	switch monolithModule {
	case 1:
		return impl.SendPaperCloudPropertyEvaluatorModuleUserVerificationEmail(ctx, email, verificationCode, firstName)
	default:
		return fmt.Errorf("unsupported monolith module: %d", monolithModule)
	}
}

func (impl *templatedEmailer) SendPaperCloudPropertyEvaluatorModuleUserVerificationEmail(ctx context.Context, email, verificationCode, firstName string) error {
	fp := path.Join("templates", "ipe/user_verification_email.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return fmt.Errorf("user verification parsing error: %w", err)
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		Email            string
		VerificationCode string
		FirstName        string
	}{
		Email:            email,
		VerificationCode: verificationCode,
		FirstName:        firstName,
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		return fmt.Errorf("user verification template execution error: %w", err)
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.incomePropertyEmailer.Send(ctx, impl.incomePropertyEmailer.GetSenderEmail(), "Activate your PaperCloud Property Evaluator Account", email, body); err != nil {
		return fmt.Errorf("sending income property evaluator user verification error: %w", err)
	}
	log.Println("success in sending income property evaluator user verification email")
	return nil
}
