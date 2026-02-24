package emails

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type Sender interface {
	Send(toEmailAddress, content string) error
}

type ResendSender struct {
	ApiKey string
}

func (sender *ResendSender) Send(toEmailAddress, content string) error {
	client := resend.NewClient(sender.ApiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{toEmailAddress},
		Subject: "Shopping List",
		Text:    content,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email in ResenderSender.Send(): %w", err)
	}
	return nil
}
