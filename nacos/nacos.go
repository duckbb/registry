package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/duckbb/registry"
)

const RegistryName = "nacos"

var Client naming_client.INamingClient

type NacosRegistry struct {
}

func (n *NacosRegistry) PluginName() string {
	return RegistryName
}

//init nacos client
func NacosInit(ClientConfig *constant.ClientConfig, ServerConfig []constant.ServerConfig) (err error) {
	Client, err = clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": ClientConfig,
		"clientConfig":  ServerConfig,
	})
	return err
}

func (n *NacosRegistry) Init(ctx context.Context, f func() error) error {
	return f()
}

//register service
func (n *NacosRegistry) Register(ctx context.Context, service *registry.Service) error {
	if Client == nil {
		return NacosNotFoundErr
	}
	resisterParam, err := NewRegisterInstanceParam(service)
	if err != nil {
		return err
	}

	success, err := Client.RegisterInstance(*resisterParam)
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
	if Client == nil {
		return NacosNotFoundErr
	}
	resisterParam := NewDeregisterInstanceParam(service)

	success, err := Client.DeregisterInstance(*resisterParam)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosDeregistionErr, err)
	}
	if !success {
		return NacosDeregistionErr
	}
	return nil
}

func (n *NacosRegistry) Get(ctx context.Context, service *registry.Service) error {
	if Client == nil {
		return NacosNotFoundErr
	}
	resisterParam := NewDeregisterInstanceParam(service)

	success, err := Client.DeregisterInstance(*resisterParam)
	if err != nil {
		return fmt.Errorf("%w,source Err:%s", NacosDeregistionErr, err)
	}
	if !success {
		return NacosDeregistionErr
	}
	return nil
}
