package nacos

import (
	"context"
	"fmt"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/nacos-group/nacos-sdk-go/clients"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/duckbb/registry"
)

const RegistryName = "nacos"

type NacosRegistry struct {
	lock     sync.Mutex
	Client   naming_client.INamingClient
	Services map[string][]*registry.Service
}

func (n *NacosRegistry) PluginName() string {
	return RegistryName
}

//init nacos client
func NacosInit(ClientConfig constant.ClientConfig, ServerConfig []constant.ServerConfig) (*NacosRegistry, error) {
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"clientConfig":  ClientConfig,
		"serverConfigs": ServerConfig,
	})
	if err != nil {
		return nil, err
	}
	r := &NacosRegistry{
		Client:   client,
		Services: map[string][]*registry.Service{},
	}
	return r, err
}

func (n *NacosRegistry) Init(ctx context.Context, f func() error) error {
	return f()
}

//register service
func (n *NacosRegistry) Register(ctx context.Context, service *registry.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	resisterParam, err := NewRegisterInstanceParam(service)
	if err != nil {
		return err
	}

	success, err := n.Client.RegisterInstance(*resisterParam)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosRegistionErr, err)
	}
	if !success {
		return NacosRegistionErr
	}
	return nil
}

//unregister service
func (n *NacosRegistry) UnRegister(ctx context.Context, service *registry.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	resisterParam := NewDeregisterInstanceParam(service)

	success, err := n.Client.DeregisterInstance(*resisterParam)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosDeregistionErr, err)
	}
	if !success {
		return NacosDeregistionErr
	}
	return nil
}

//get services
func (n *NacosRegistry) Get(ctx context.Context, service *registry.Service) ([]*registry.Service, error) {
	if srvs, ok := n.Services[service.NacosServiceName]; ok {
		return srvs, nil
	}
	if n.Client == nil {
		return nil, NacosNotFoundErr
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	//service healthy=true,enable=true 和weight>0
	param := NewSelectInstances(service)
	instances, err := n.Client.SelectInstances(*param)
	if err != nil {
		return nil, fmt.Errorf("%w,source Err:%s", NacosGetServiceErr, err)
	}
	srvs := []*registry.Service{}
	for _, v := range instances {
		srv := &registry.Service{
			NacosIp:          v.Ip,
			NacosPort:        v.Port,
			NacosWeight:      v.Weight,
			NacosHealthy:     v.Healthy,
			NacosMetadata:    v.Metadata,
			NacosClusterName: v.ClusterName,
			NacosServiceName: v.ServiceName,
			NacosEphemeral:   false,
		}
		srvs = append(srvs, srv)
	}
	n.Services[service.NacosServiceName] = srvs
	return srvs, nil

}

func NewNacosService(ServiceName, Ip string, Port uint64, Weight float64, Enable, Healthy bool, opts ...NacosServiceOption) *registry.Service {
	srv := &registry.Service{
		NacosIp:          Ip,
		NacosPort:        Port,
		NacosWeight:      Weight,
		NacosEnable:      Enable,
		NacosHealthy:     Healthy,
		NacosServiceName: ServiceName,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func NewRegisterInstanceParam(srv *registry.Service) (*vo.RegisterInstanceParam, error) {
	c := &vo.RegisterInstanceParam{
		Ip:          srv.NacosIp,
		Port:        srv.NacosPort,
		Weight:      srv.NacosWeight,
		Enable:      srv.NacosEnable,
		Healthy:     srv.NacosHealthy,
		Metadata:    srv.NacosMetadata,
		ClusterName: srv.NacosClusterName,
		ServiceName: srv.NacosServiceName,
		GroupName:   srv.NacosGroupName,
		Ephemeral:   srv.NacosEphemeral,
	}
	if c.Weight <= 0 {
		return nil, NacosWeightErr
	}
	return c, nil
}

func NewDeregisterInstanceParam(srv *registry.Service) *vo.DeregisterInstanceParam {
	c := &vo.DeregisterInstanceParam{
		Ip:          srv.NacosIp,
		Port:        srv.NacosPort,
		Cluster:     srv.NacosClusterName,
		ServiceName: srv.NacosServiceName,
		GroupName:   srv.NacosGroupName,
		Ephemeral:   srv.NacosEphemeral,
	}
	return c
}

func NewSelectInstances(srv *registry.Service) *vo.SelectInstancesParam {
	c := &vo.SelectInstancesParam{
		ServiceName: srv.NacosServiceName,
		HealthyOnly: true,
	}
	if srv.NacosClusterName != "" {
		c.Clusters = []string{srv.NacosClusterName}
	}
	if srv.NacosGroupName != "" {
		c.GroupName = srv.NacosGroupName
	}
	return c
}
