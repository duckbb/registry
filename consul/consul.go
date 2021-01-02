package consul

import (
	"sync"

	vo2 "github.com/duckbb/registry/base"
	"github.com/hashicorp/consul/api"
)

const RegistryName = "consul"

type ConsulRegistry struct {
	sync.Mutex
	Client   *api.Client
	Services map[string][]*vo2.Service
}

func (n *ConsulRegistry) PluginName() string {
	return RegistryName
}

//init consul client
func NacosInit(c *api.Config) (*ConsulRegistry, error) {
	client, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}
	return &ConsulRegistry{
		Client:   client,
		Services: map[string][]*vo2.Service{},
	}, err
}
