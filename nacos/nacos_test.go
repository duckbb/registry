package nacos

import (
	"context"
	"math/rand"
	"strconv"
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

//register same service
func TestEqualService(t *testing.T) {
	c, err := getNacosRegistry()
	if err != nil {
		t.Errorf("get client fail,err:%s", err)
	}
	service := &registry.Service{
		NacosServiceName: "demo.go",
		NacosIp:          "192.168.50.229",
		NacosPort:        8090,
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

	service1 := &registry.Service{
		NacosServiceName: "demo.go",
		NacosIp:          "192.168.50.229",
		NacosPort:        8090,
		NacosWeight:      20,
		NacosEnable:      true,
		NacosHealthy:     true,
		NacosMetadata:    map[string]string{"idc": "shanghai11111"},
		NacosEphemeral:   true,
	}
	//overlay
	err = c.Register(context.TODO(), service1)
	if err != nil {
		t.Errorf("equal register service fail,err:%s", err)
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

//test subscribe case
func TestSubscribe(t *testing.T) {
	c, err := getNacosRegistry()
	if err != nil {
		t.Errorf("get client fail,err:%s", err)
	}
	//register
	go func() {
		c2, err := getNacosRegistry()
		if err != nil {
			t.Errorf("get client fail,err:%s", err)
		}
		serviceNum := 2
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < serviceNum; i++ {
			ser := &registry.Service{
				NacosServiceName: "demo.go",
				NacosIp:          "192.168.50.229",
				NacosPort:        uint64(rand.Intn(10000)),
				NacosWeight:      20,
				NacosEnable:      true,
				NacosHealthy:     true,
				NacosEphemeral:   true,
			}
			ser.NacosMetadata = map[string]string{"idc": "shanghai" + strconv.FormatUint(ser.NacosPort, 10)}
			err := c2.Register(context.TODO(), ser)
			if err != nil {
				t.Log("register service fail,err", err.Error())
			}
			t.Logf("register one sevice:servicePort:%d", ser.NacosPort)
			time.Sleep(time.Second * 3)
		}
	}()
	//watch
	service := &registry.Service{
		NacosServiceName: "demo.go",
	}
	err = c.SubscribeService(context.TODO(), service)
	if err != nil {
		t.Errorf("subscribe service failed,err:%s", err)
	}
	go func() {
		for {
			t.Log("----print service data start----")
			if services, ok := c.Services[service.NacosServiceName]; ok {
				for _, v := range services {
					t.Logf("service data:%+v", v)
				}
			}
			t.Log("----print service data stop----")

			time.Sleep(time.Second * 5)
		}

	}()
	select {}
	//time.Sleep(time.Second * 20)
}
