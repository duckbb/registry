package main

/*
http://192.168.50.78:8500/
*/
import (
	"context"
	"fmt"
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	consulAddress = "http://192.168.50.78:8500"
	localIp       = "192.168.50.229"
	localPort     = 81
)

func consulRegister(port int) error {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "337"
	registration.Name = "service337"
	registration.Port = localPort
	registration.Tags = []string{"testService"}
	registration.Address = localIp

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	chanErr := make(chan error)
	ports := []int{8085, 8086}
	for _, port := range ports {
		go func() {
			srv1 := NewServer(localIp + ":" + strconv.Itoa(port))
			err := consulRegister(port)
			if err != nil {
				chanErr <- err
				return
			}
			err = srv1.Start(context.Background())
			if err != nil {
				chanErr <- err
				return
			}
		}()
	}

	select {
	case err := <-chanErr:
		log.Println("err:", err)
	}
}
