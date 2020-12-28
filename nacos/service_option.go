package nacos

import (
	"github.com/duckbb/registry"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

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

type NacosServiceOption func(*registry.Service)

func WithMetadata(data map[string]string) NacosServiceOption {
	return func(service *registry.Service) {
		service.NacosMetadata = data
	}
}

func WithClusterName(ClusterName string) NacosServiceOption {
	return func(service *registry.Service) {
		service.NacosClusterName = ClusterName
	}
}

func WithGroupName(GroupName string) NacosServiceOption {
	return func(service *registry.Service) {
		service.NacosGroupName = GroupName
	}
}

func WithEphemeral(Ephemeral bool) NacosServiceOption {
	return func(service *registry.Service) {
		service.NacosEphemeral = Ephemeral
	}
}
