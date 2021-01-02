package consul

import (
	"context"
	"fmt"
	"strconv"
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

	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = "337_" + strconv.Itoa(service.ConsulPort)
	registration.Name = "service337_"
	registration.Port = port
	registration.Tags = []string{"testService"}
	registration.Address = localIp

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "3s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "5s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

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
