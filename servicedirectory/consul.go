package servicedirectory

import (
	"context"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/soupstoregames/go-core/logging"
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

func (c *ConsulServiceCatalog) Service(ctx context.Context, service string, tag string) ([]string, error) {
	for {
		services, _, err := c.catalog.Service(service, tag, nil)
		if err != nil {
			logging.Error(err.Error())
		}

		select {
		case <-ctx.Done():
			logging.Error(fmt.Sprintf("Consul: Finding service '%s' tagged '%s': %v", service, tag, ctx.Err()))
			return []string{}, ctx.Err()
		default:
			if len(services) == 0 {
				continue
			}

			addresses := make([]string, len(services))
			for i := range services {
				addresses[i] = fmt.Sprintf("%s:%d", services[i].ServiceAddress, services[i].ServicePort)
			}

			return addresses, nil
		}
	}
}
