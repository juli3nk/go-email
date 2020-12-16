package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	data, err := json.Marshal(req)
        if err != nil {
                return err
        }

        fakemailUrl := fmt.Sprintf("%s/incoming", p.Config["url"])
        contentType := "application/json"

        tr := &http.Transport{
                IdleConnTimeout: 5 * time.Second,
        }
        client := &http.Client{Transport: tr}

        response, err := client.Post(fakemailUrl, contentType, bytes.NewBuffer(data))
        if err != nil {
                return err
        }

        if response.StatusCode != 201 {
                return fmt.Errorf("Please contact administrator if the problem persists")
        }

	return nil
}
