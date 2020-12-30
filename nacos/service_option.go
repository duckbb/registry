package nacos

import (
	vo "github.com/duckbb/registry/base"
)

type NacosServiceOption func(*vo.Service)

func WithMetadata(data map[string]string) NacosServiceOption {
	return func(service *vo.Service) {
		service.NacosMetadata = data
	}
}

func WithClusterName(ClusterName string) NacosServiceOption {
	return func(service *vo.Service) {
		service.NacosClusterName = ClusterName
	}
}

func WithGroupName(GroupName string) NacosServiceOption {
	return func(service *vo.Service) {
		service.NacosGroupName = GroupName
	}
}

func WithEphemeral(Ephemeral bool) NacosServiceOption {
	return func(service *vo.Service) {
		service.NacosEphemeral = Ephemeral
	}
}
