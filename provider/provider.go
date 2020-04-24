package provider

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/juli3nk/go-email/types"
)

type Provider interface {
	Send(*types.SendRequest) error
}

type ProviderFactory func(conf map[string]string) (Provider, error)

var providerFactories = make(map[string]ProviderFactory)

func supportedDrivers() string {
	drivers := make([]string, 0, len(providerFactories))

	for d := range providerFactories {
		drivers = append(drivers, string(d))
	}

	sort.Strings(drivers)

	return strings.Join(drivers, ",")
}

func RegisterDriver(name string, factory ProviderFactory) {
	if factory == nil {
		log.Panicf("Provider factory %s does not exist.", name)
	}

	if _, registered := providerFactories[name]; registered {
		log.Printf("Provider factory %s already registered. Ignoring.", name)
	}

	providerFactories[name] = factory
}

func NewDriver(driver string, config map[string]string) (Provider, error) {
	engineFactory, exists := providerFactories[driver]
	if exists {
		return engineFactory(config)
	}

	return nil, fmt.Errorf("The driver: %s is not supported. Supported drivers are %s", driver, supportedDrivers())
}
