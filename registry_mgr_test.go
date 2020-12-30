package registry

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

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
	plugin.SetRegisterType(NACOS)
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

func TestGetService(t *testing.T) {
	c, err := getMgr()
	if err != nil {
		t.Errorf("get mgr fail,err:%s\n", err)
	}
	service := &vo.Service{
		NacosServiceName: "demo.go",
	}
	services, err := c.Get(context.TODO(), service)
	if err != nil {
		t.Errorf("get base failed,err:%s", err)
	}
	for _, v := range services {
		t.Logf("%+v", v)
	}

}

//test subscribe case
func TestSubscribe(t *testing.T) {
	c, err := getMgr()
	if err != nil {
		t.Errorf("get mgr fail,err:%s\n", err)
	}
	//register
	go func() {
		if err != nil {
			t.Errorf("get client fail,err:%s", err)
		}
		serviceNum := 2
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < serviceNum; i++ {
			ser := &vo.Service{
				NacosServiceName: "demo.go",
				NacosIp:          "192.168.50.229",
				NacosPort:        uint64(rand.Intn(10000)),
				NacosWeight:      20,
				NacosEnable:      true,
				NacosHealthy:     true,
				NacosEphemeral:   true,
			}
			ser.NacosMetadata = map[string]string{"idc": "shanghai" + strconv.FormatUint(ser.NacosPort, 10)}
			err := c.Register(context.TODO(), ser)
			if err != nil {
				t.Log("register base fail,err", err.Error())
			}
			t.Logf("register one sevice:servicePort:%d", ser.NacosPort)
			time.Sleep(time.Second * 3)
		}
	}()
	//watch
	service := &vo.Service{
		NacosServiceName: "demo.go",
	}
	err = c.SubscribeService(context.TODO(), service)
	if err != nil {
		t.Errorf("subscribe base failed,err:%s", err)
	}
	clent, err := c.GetRegister()
	if err != nil {
		t.Error("get client fail")
	}
	nacosClient := clent.(*nacos.NacosRegistry)
	go func() {
		for {
			t.Log("----print base data start----")

			if services, ok := nacosClient.Services[service.NacosServiceName]; ok {
				for _, v := range services {
					t.Logf("base data:%+v", v)
				}
			}
			t.Log("----print base data stop----")

			time.Sleep(time.Second * 5)
		}

	}()
	select {}
	//time.Sleep(time.Second * 50)
}
