package nacos

import (
	"context"
	"fmt"
	"log"
	"sync"

	vo2 "github.com/duckbb/registry/base"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const RegistryName = "nacos"

type NacosRegistry struct {
	sync.Mutex
	Client   naming_client.INamingClient
	Services map[string][]*vo2.Service
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
		Services: map[string][]*vo2.Service{},
	}
	return r, err
}

func (n *NacosRegistry) Init(ctx context.Context, f func() error) error {
	return f()
}

//register base
func (n *NacosRegistry) Register(ctx context.Context, service *vo2.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	n.Lock()
	defer n.Unlock()

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

//unregister base
func (n *NacosRegistry) UnRegister(ctx context.Context, service *vo2.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	n.Lock()
	defer n.Unlock()
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
func (n *NacosRegistry) Get(ctx context.Context, service *vo2.Service) ([]*vo2.Service, error) {
	if srvs, ok := n.Services[service.NacosServiceName]; ok {
		return srvs, nil
	}
	if n.Client == nil {
		return nil, NacosNotFoundErr
	}
	n.Lock()
	defer n.Unlock()
	//base healthy=true,enable=true 和weight>0
	param := NewSelectInstances(service)
	instances, err := n.Client.SelectInstances(*param)
	if err != nil {
		return nil, fmt.Errorf("%w,source Err:%s", NacosGetServiceErr, err)
	}
	srvs := []*vo2.Service{}
	for _, v := range instances {
		srv := &vo2.Service{
			NacosIp:          v.Ip,
			NacosPort:        v.Port,
			NacosWeight:      v.Weight,
			NacosEnable:      v.Enable,
			NacosHealthy:     v.Healthy,
			NacosMetadata:    v.Metadata,
			NacosServiceName: v.ServiceName,
			NacosClusterName: v.ClusterName,
			NacosEphemeral:   v.Ephemeral,
		}
		if service.NacosGroupName != "" {
			srv.NacosGroupName = service.NacosGroupName
		}
		srvs = append(srvs, srv)
	}
	n.Services[service.NacosServiceName] = srvs
	return srvs, nil

}

//subscribe Service
func (n *NacosRegistry) SubscribeService(ctx context.Context, service *vo2.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	//param := NewSubscribeParam(base)
	param := &vo.SubscribeParam{
		ServiceName: service.NacosServiceName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Println("watch base change:")
			srvs := []*vo2.Service{}
			for _, v := range services {
				tempService := &vo2.Service{
					NacosIp:          v.Ip,
					NacosPort:        v.Port,
					NacosWeight:      v.Weight,
					NacosEnable:      v.Enable,
					NacosMetadata:    v.Metadata,
					NacosClusterName: v.ClusterName,
					NacosServiceName: v.ServiceName,
				}
				srvs = append(srvs, tempService)
				log.Printf("watch base:%+v\n", v)
			}
			n.Lock()
			defer n.Unlock()
			n.Services[service.NacosServiceName] = srvs
		},
	}
	if service.NacosGroupName != "" {
		param.GroupName = service.NacosGroupName
	}
	if service.NacosClusterName != "" {
		param.Clusters = []string{service.NacosClusterName}
	}
	err := n.Client.Subscribe(param)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosSubscribeErr, err)
	}
	return nil
}

//unsubscribe base
func (n *NacosRegistry) UnsubscribeService(ctx context.Context, service *vo2.Service) error {
	if n.Client == nil {
		return NacosNotFoundErr
	}
	param := &vo.SubscribeParam{
		ServiceName: service.NacosServiceName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Printf("unwatch base:%+v\n", services)
		},
	}
	if service.NacosGroupName != "" {
		param.GroupName = service.NacosGroupName
	}
	if service.NacosClusterName != "" {
		param.Clusters = []string{service.NacosClusterName}
	}
	err := n.Client.Unsubscribe(param)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosUnsubscribeErr, err)
	}
	return nil
}

//equal base
func EqualService(s1, s2 *vo2.Service) bool {
	if s1.NacosServiceName == s2.NacosServiceName &&
		s1.NacosIp == s2.NacosIp &&
		s1.NacosPort == s2.NacosPort &&
		s1.NacosGroupName == s2.NacosGroupName &&
		s1.NacosClusterName == s2.NacosClusterName {
		return true
	}
	return false
}

func NewNacosService(ServiceName, Ip string, Port uint64, Weight float64, Enable, Healthy bool, opts ...NacosServiceOption) *vo2.Service {
	srv := &vo2.Service{
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

func NewRegisterInstanceParam(srv *vo2.Service) (*vo.RegisterInstanceParam, error) {
	param := &vo.RegisterInstanceParam{
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
	if param.Weight <= 0 {
		return nil, NacosWeightErr
	}
	return param, nil
}

func NewDeregisterInstanceParam(srv *vo2.Service) *vo.DeregisterInstanceParam {
	param := &vo.DeregisterInstanceParam{
		Ip:          srv.NacosIp,
		Port:        srv.NacosPort,
		Cluster:     srv.NacosClusterName,
		ServiceName: srv.NacosServiceName,
		GroupName:   srv.NacosGroupName,
		Ephemeral:   srv.NacosEphemeral,
	}
	return param
}

func NewSelectInstances(srv *vo2.Service) *vo.SelectInstancesParam {
	param := &vo.SelectInstancesParam{
		ServiceName: srv.NacosServiceName,
		HealthyOnly: true,
	}
	if srv.NacosClusterName != "" {
		param.Clusters = []string{srv.NacosClusterName}
	}
	if srv.NacosGroupName != "" {
		param.GroupName = srv.NacosGroupName
	}
	return param
}
