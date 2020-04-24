package provider

import (
	"fmt"

	"github.com/go-resty/resty"
	"github.com/juli3nk/go-email/types"
)

func init() {
	RegisterDriver("fakemail", NewFakeMailProvider)
}

type FakeMailProvider struct {
	Config map[string]string
}

func NewFakeMailProvider(config map[string]string) (Provider, error) {
	if _, ok := config["url"]; !ok {
		config["url"] = "http://fakemail.local:8080"
	}

	return &FakeMailProvider{Config: config}, nil
}

func (p *FakeMailProvider) Send(req *types.SendRequest) error {
	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(req).
		Post(fmt.Sprintf("%s/incoming", p.Config["url"]))
	if err != nil {
		return err
	}

	return nil
}
