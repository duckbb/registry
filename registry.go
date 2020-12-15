package registry

import "context"

//plunin interface
type Registryer interface {
	Name() string
	Init(context.Context) error
	Register(context.Context, *Service) error
	UnRegister(context.Context, *Service) error
	Get(context.Context, *Service) error
}

//service
type Service struct {
	Name  string  `json:"name"`
	nodes []*Node `json:"nodes"`
}

// node
type Node struct {
	Addr   string `json:"addr"`
	Group  string `json:"group"` //label for group Service
	Region string `json:"region"`
	Zone   string `json:"zone"`

	Weight int `json:"weight"` // weight round roubin
}
