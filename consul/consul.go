package consul

import (
	"context"
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

func (n *ConsulRegistry) Init(ctx context.Context, f func() error) error {
	return f()
}

//register base
func (n *ConsulRegistry) Register(ctx context.Context, service *vo2.Service) error {
	if n.Client == nil {
		return ConsulNotFoundErr
	}
	n.Lock()
	defer n.Unlock()

	//resisterParam, err := NewRegisterInstanceParam(service)
	//if err != nil {
	//	return err
	//}
	//
	//success, err := n.Client.RegisterInstance(*resisterParam)
	//if err != nil {
	//	return fmt.Errorf("%w,source Err:%s", NacosRegistionErr, err)
	//}
	//if !success {
	//	return NacosRegistionErr
	//}
	//return nil
}
