package registry

import "context"

//plunin interface
type Registryer interface {
	Name() string
	Init(context.Context, func() error) error
	Register(context.Context, *Service) error
	UnRegister(context.Context, *Service) error
	Get(context.Context, *Service) ([]*Service, error)
	SubscribeService(context.Context, *Service) error
}

type Service struct {

	//nacos init
	NacosIp          string
	NacosPort        uint64
	NacosWeight      float64
	NacosEnable      bool
	NacosHealthy     bool
	NacosMetadata    map[string]string
	NacosClusterName string
	NacosServiceName string
	NacosGroupName   string
	NacosEphemeral   bool
}

////service
//type Service struct {
//	Name  string  `json:"name"`
//	nodes []*Node `json:"nodes"`
//}
//
//// node
//type Node struct {
//	Ip       string `json:"ip"`
//	Port     string `json:"port"`
//	Group    string `json:"group"` //label for group Service
//	Region   string `json:"region"`
//	Zone     string `json:"zone"`
//	Metadata map[string]string
//	Weight   int `json:"weight"` // weight round roubin
//
//	Ip:          "10.0.0.11",
//	Port:        8848,
//	ServiceName: "demo.go",
//	Weight:      10,
//	Enable:      true,
//	Healthy:     true,
//	Ephemeral:   true,
//	Metadata:    map[string]string{"idc":"shanghai"},
//	ClusterName: "cluster-a", // 默认值DEFAULT
//	GroupName:   "group-a",  // 默认值DEFAULT_GROUP
//}
