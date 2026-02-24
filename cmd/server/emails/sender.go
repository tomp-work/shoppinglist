package emails

import (
	"github.com/resend/resend-go/v2"
)

type Sender interface {
	Send(toEmail, content string) error
}

type ResendSender struct {
	ApiKey string
}

func (sender *ResendSender) Send(toEmail, content string) error {
	client := resend.NewClient(sender.ApiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{toEmail},
		Subject: "Shopping List",
		Text:    content,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}
