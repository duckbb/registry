package nacos

import (
	"github.com/duckbb/registry"
)

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
