package consul

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

type ConsulWatcher struct {
	resolver *ConsulResolver
	client   *consul.Client

	lastIndex uint64
	addresses []string
}

func (w *ConsulWatcher) Next() ([]*naming.Update, error) {
	if w.addresses == nil {
		addresses, lastIndex, _ := w.queryConsul(nil)
		if len(addresses) != 0 {
			w.addresses = addresses
			w.lastIndex = lastIndex
			return GenUpdates([]string{}, addresses), nil
		}
	}

	for {
		addresses, lastIndex, err := w.queryConsul(&consul.QueryOptions{WaitIndex: w.lastIndex})
		if err != nil {
			continue
		}
		updates := GenUpdates(w.addresses, addresses)
		w.addresses = addresses
		w.lastIndex = lastIndex

		if len(updates) != 0 {
			return updates, nil
		}
	}
	return []*naming.Update{}, nil
}

func (w *ConsulWatcher) Close() {
}

func (w *ConsulWatcher) queryConsul(q *consul.QueryOptions) ([]string, uint64, error) {
	cs, meta, err := w.client.Health().Service(w.resolver.ServiceName, "", true, q)
	if err != nil {
		return nil, 0, err
	}

	addrs := make([]string, len(cs))
	for i, s := range cs {
		addrs[i] = fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port)
	}
	return addrs, meta.LastIndex, nil
}
