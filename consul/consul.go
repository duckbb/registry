package consul

import (
	"context"
	"fmt"
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

func (c *ConsulRegistry) PluginName() string {
	return RegistryName
}

//init consul client
func ConsulInit(c *api.Config) (*ConsulRegistry, error) {
	client, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}
	return &ConsulRegistry{
		Client:   client,
		Services: map[string][]*vo2.Service{},
	}, err
}

func (c *ConsulRegistry) Init(ctx context.Context, f func() error) error {
	return f()
}

//register base
func (c *ConsulRegistry) Register(ctx context.Context, service *vo2.Service) error {
	if c.Client == nil {
		return ConsulNotFoundErr
	}
	c.Lock()
	defer c.Unlock()

	AgentService := NewAgentServiceRegistration(service)

	err := c.Client.Agent().ServiceRegister(AgentService)
	if err != nil {
		return fmt.Errorf("%w,register fail,err:%s", ConsulRegistionErr, err)
	}
	return nil
}

//create agentService
func NewAgentServiceRegistration(service *vo2.Service) *api.AgentServiceRegistration {
	agent := &api.AgentServiceRegistration{
		Kind:              service.ConsulKind,
		ID:                service.ConsulID,
		Name:              service.ConsulName,
		Tags:              service.ConsulTags,
		Port:              service.ConsulPort,
		Address:           service.ConsulAddress,
		TaggedAddresses:   service.ConsulTaggedAddresses,
		EnableTagOverride: service.ConsulEnableTagOverride,
		Meta:              service.ConsulMeta,
		Weights:           service.ConsulWeights,
		Check:             service.ConsulCheck,
		Checks:            service.ConsulChecks,
		Proxy:             service.ConsulProxy,
		Connect:           service.ConsulConnect,
		Namespace:         service.ConsulNamespace,
	}
	return agent
}
