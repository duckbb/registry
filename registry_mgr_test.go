package registry

import (
	"context"
	"log"
	"testing"

	vo "github.com/duckbb/registry/base"

	"github.com/duckbb/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func TestGetRegistryMgr(t *testing.T) {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.50.78",
			Port:        8848,
			Scheme:      "http",
			ContextPath: "/nacos",
		},
	}
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		RotateTime:          "1h",
		MaxAge:              3,
		Username:            "nacos",
		Password:            "123456",
	}
	registry, err := nacos.NacosInit(clientConfig, serverConfigs)
	if err != nil {
		t.Errorf("init nacos registry fail,err:%s\n", err)
	}
	plugin := GetRegistryMgr()
	err = plugin.InitPluninRegistry(context.TODO(), NACOS, registry)
	if err != nil {
		t.Errorf("plunin registry fail,err:%s\n", err)
	}

}

//get plugin

func getMgr() (*pluninRegistry, error) {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.50.78",
			Port:        8848,
			Scheme:      "http",
			ContextPath: "/nacos",
		},
	}
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		RotateTime:          "1h",
		MaxAge:              3,
		Username:            "nacos",
		Password:            "123456",
	}
	registry, err := nacos.NacosInit(clientConfig, serverConfigs)
	if err != nil {
		return nil, err
	}
	plugin := GetRegistryMgr()
	//plugin.SetRegisterType(NACOS)
	err = plugin.InitPluninRegistry(context.TODO(), NACOS, registry)
	if err != nil {
		return nil, err
	}
	return plugin, nil
}

//register service
func TestRegister(t *testing.T) {
	mgr, err := getMgr()
	if err != nil {
		t.Errorf("get mgr fail,err:%s\n", err)
	}
	log.Println("err:", err)
	service := &vo.Service{
		NacosServiceName: "demo.go",
		NacosIp:          "192.168.50.229",
		NacosPort:        8081,
		NacosWeight:      10,
		NacosEnable:      true,
		NacosHealthy:     true,
		NacosMetadata:    map[string]string{"idc": "shanghai11111"},
		NacosEphemeral:   true,
	}
	err = mgr.Register(context.TODO(), service, NACOS)
	if err != nil {
		t.Errorf("register base fail,err:%s", err)
	}

}
