package servicedirectory

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

type ConsulServiceCatalog struct {
	catalog *consul.Catalog
}

func NewConsulServiceCatalog() (*ConsulServiceCatalog, error) {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}

	catalog := consulClient.Catalog()

	return &ConsulServiceCatalog{catalog}, nil
}

func (c *ConsulServiceCatalog) Service(service string) ([]string, error) {
	services, _, err := c.catalog.Service(service, "", nil)
	if err != nil {
		return []string{}, err
	}

	addresses := make([]string, len(services))
	for i := range services {
		addresses[i] = fmt.Sprintf("%s:%d", services[i].ServiceAddress, services[i].ServicePort)
	}

	return addresses, nil
}
