package nacos

import (
	"testing"
	"time"

	"github.com/duckbb/registry"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func getClient() (*NacosRegistry, error) {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.50.78",
			Port:        8848,
			Scheme:      "http",
			ContextPath: "/nacos",
		},
	}
	clientConfig := constant.ClientConfig{
		//NamespaceId:         "82934ebb-070f-4193-b843-2053b3782563", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//LogDir:              "E:\\go_learning\\a4\\Go-000\\Week04\\cmd\\article\\log",
		//CacheDir:            "E:\\go_learning\\a4\\Go-000\\Week04\\cmd\\article\\cache",
		//LogDir:     "/tmp/nacos/log",
		//CacheDir:   "/tmp/nacos/cache",
		RotateTime: "1h",
		MaxAge:     3,
		//LogLevel:   "debug",
		Username: "nacos",
		Password: "123456",
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
		//NamespaceId:         "82934ebb-070f-4193-b843-2053b3782563", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//LogDir:              "E:\\go_learning\\a4\\Go-000\\Week04\\cmd\\article\\log",
		//CacheDir:            "E:\\go_learning\\a4\\Go-000\\Week04\\cmd\\article\\cache",
		//LogDir:     "/tmp/nacos/log",
		//CacheDir:   "/tmp/nacos/cache",
		RotateTime: "1h",
		MaxAge:     3,
		//LogLevel:   "debug",
		Username: "nacos",
		Password: "123456",
	}
	_, err := NacosInit(clientConfig, serverConfigs)
	if err != nil {
		t.Errorf("nacos init fail:%s", err)
	}

}

func TestNacosRegister(t *testing.T) {
	c, err := getClient()
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
	v1, err := NewRegisterInstanceParam(service)
	if err != nil {
		t.Error(err)
	}
	success, err := c.Client.RegisterInstance(*v1)
	if err != nil {
		t.Errorf("register service fail,err:%s", err)
	}
	if !success {
		t.Errorf("final register fail")
	}
	time.Sleep(time.Second * 10)
}
