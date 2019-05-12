package grpc

import (
	"context"
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

	conn, err := grpc.Dial(chunkAddr[0], grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
