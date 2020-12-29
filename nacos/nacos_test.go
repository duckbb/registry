package nacos

import (
	"context"
	"testing"
	"time"

	"github.com/duckbb/registry"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func getNacosRegistry() (*NacosRegistry, error) {
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
	registry, err := NacosInit(clientConfig, serverConfigs)
	return registry, err
}

func TestNacosInit(t *testing.T) {
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
	_, err := NacosInit(clientConfig, serverConfigs)
	if err != nil {
		t.Errorf("nacos init fail:%s", err)
	}

}

func TestNacosRegister(t *testing.T) {
	c, err := getNacosRegistry()
	if err != nil {
		t.Errorf("get client fail,err:%s", err)
	}
	service := &registry.Service{
		NacosServiceName: "demo.go",
		NacosIp:          "192.168.50.229",
		NacosPort:        8081,
		NacosWeight:      10,
		NacosEnable:      true,
		NacosHealthy:     true,
		NacosMetadata:    map[string]string{"idc": "shanghai11111"},
		NacosEphemeral:   true,
	}
	err = c.Register(context.TODO(), service)
	if err != nil {
		t.Errorf("register service fail,err:%s", err)
	}
	time.Sleep(time.Second * 10)
}

func TestGetService(t *testing.T) {
	c, err := getNacosRegistry()
	if err != nil {
		t.Errorf("get client fail,err:%s", err)
	}
	service := &registry.Service{
		NacosServiceName: "demo.go",
	}
	services, err := c.Get(context.TODO(), service)
	if err != nil {
		t.Errorf("get service failed,err:%s", err)
	}
	for _, v := range services {
		t.Logf("%+v", v)
	}

}
