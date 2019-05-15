package grpc

import (
	"context"
	"fmt"
	"github.com/soupstoregames/go-core/logging"
	"github.com/soupstoregames/go-core/servicedirectory"
	"google.golang.org/grpc"
)

type ConsulDialler struct {
	catalog *servicedirectory.ConsulServiceCatalog
}

func NewConsulDialler(catalog *servicedirectory.ConsulServiceCatalog) *ConsulDialler {
	return &ConsulDialler{
		catalog: catalog,
	}
}

func (d *ConsulDialler) Dial(ctx context.Context, service, tag string) (*grpc.ClientConn, error) {
	chunkAddr, err := d.catalog.Service(ctx, service, tag)
	if err != nil {
		return nil, err
	}

	for {
		conn, err := grpc.Dial(chunkAddr[0], grpc.WithInsecure())
		if err != nil {
			logging.Info("Failed to connect to %s tagged %s, retrying...")
		}

		select {
		case <-ctx.Done():
			logging.Error(fmt.Sprintf("Consul: Finding service '%s' tagged '%s': %v", service, tag, ctx.Err()))
			return nil, ctx.Err()
		default:
			if err != nil {
				continue
			}

			return conn, nil
		}
	}

}
