package consul

import (
	"errors"
	"fmt"

	consul "github.com/hashicorp/consul/api"

	"google.golang.org/grpc/naming"
)

type ConsulResolver struct {
	ServiceName string
}

func (r *ConsulResolver) Resolve(target string) (naming.Watcher, error) {
	if r.ServiceName == "" {
		return nil, errors.New("can't resolve with empty service name")
	}

	conf := &consul.Config{
		Address: target,
	}
	client, err := consul.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("can't create consul client: address = %s, err = %v", target, err)
	}

	watcher := &ConsulWatcher{
		resolver: r,
		client:   client,
	}
	return watcher, nil
}

func NewResolver(serviceName string) *ConsulResolver {
	return &ConsulResolver{ServiceName: serviceName}
}
