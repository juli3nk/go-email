package provider

import (
	"fmt"

	"github.com/juli3nk/go-email/types"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func init() {
	RegisterDriver("sendgrid", NewSendGridProvider)
}

type SendGridProvider struct {
	Config map[string]string
}

func NewSendGridProvider(config map[string]string) (Provider, error) {
	if _, ok := config["api-key"]; !ok {
		return nil, fmt.Errorf("\"api-key\" config key is mandatory")
	}

	return &SendGridProvider{Config: config}, nil
}

func (p *SendGridProvider) Send(req *types.SendRequest) error {
	from := mail.NewEmail(req.From.Name, req.From.Address)

	contentText := mail.NewContent("text/plain", req.Body.Text)
	contentHtml := mail.NewContent("text/html", req.Body.Html)

	mp := mail.NewPersonalization()

	for _, a := range req.ToAddresses {
		mp.AddTos(mail.NewEmail("", a))
	}

	if len(req.CcAddresses) > 0 {
		for _, a := range req.CcAddresses {
			mp.AddCCs(mail.NewEmail("", a))
		}
	}

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = req.Subject
	message.AddPersonalizations(mp)
	message.AddContent(contentText)
	message.AddContent(contentHtml)

	if len(req.From.ReplyTo) > 0 {
		message.SetReplyTo(mail.NewEmail(req.From.Name, req.From.ReplyTo))
	}

	client := sendgrid.NewSendClient(p.Config["api-key"])

	if _, err := client.Send(message); err != nil {
		return err
	}

	return nil
}
