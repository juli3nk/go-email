package email

import (
	"fmt"

	"github.com/juli3nk/go-email/provider"
	"github.com/juli3nk/go-email/types"
)

type Config struct {
	provider provider.Provider
	request  *types.SendRequest
}

func New(providerName string, options map[string]string) (*Config, error) {
	p, err := provider.NewDriver(providerName, options)
	if err != nil {
		return nil, err
	}

	return &Config{
		provider: p,
	}, nil
}

func (c *Config) SetSender(name, address, replyToAddress string) error {
	from := types.From{
		Name:    name,
		Address: address,
	}

	if len(replyToAddress) > 0 {
		from.ReplyTo = replyToAddress
	}

	req := types.SendRequest{
		From: &from,
	}

	c.request = &req

	return nil
}

func (c *Config) SetRecipient(addresses []string) error {
	c.request.ToAddresses = addresses

	return nil
}

func (c *Config) AddCc(addresses []string) error {
	c.request.CcAddresses = addresses

	return nil
}

func (c *Config) SetSubject(value string) error {
	c.request.Subject = value

	return nil
}

func (c *Config) SetBody(bodyText, bodyHtml string) error {
	body := types.Body{
		Text: bodyText,
		Html: bodyHtml,
	}

	c.request.Body = &body

	return nil
}

func (c *Config) Send() error {
	if len(c.request.From.Name) == 0 {
		return fmt.Errorf("Sender name field is mandatory")
	}
	if len(c.request.From.Address) == 0 {
		return fmt.Errorf("Sender address field is mandatory")
	}

	if len(c.request.ToAddresses) == 0 {
		return fmt.Errorf("Recipient address field is mandatory")
	}

	if len(c.request.Body.Text) == 0 || len(c.request.Body.Html) == 0 {
		return fmt.Errorf("Body text and / or html field is mandatory")
	}

	if err := c.provider.Send(c.request); err != nil {
		return err
	}

	return nil
}
