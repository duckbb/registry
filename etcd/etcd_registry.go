package etcd

import (
	"context"
	"fmt"
	"sync"

	"github.com/coreos/etcd/clientv3"

	"github.com/duckbb/registry"
)

type EtcdRegistry struct {
	sync.Mutex
	Client *clientv3.Client
	config *clientv3.Config
}

func NewEtcdOptions(Endpoints []string, opts ...Options) *clientv3.Config {
	c := &clientv3.Config{
		Endpoints: Endpoints,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func NewEtcdRegistry(Endpoints []string, opts ...Options) (*EtcdRegistry, error) {
	c := &clientv3.Config{
		Endpoints: Endpoints,
	}
	for _, opt := range opts {
		opt(c)
	}
	cli, err := clientv3.New(*c)
	if err != nil {
		err = fmt.Errorf("etcd connect failed, err:%v", err)
		return nil, err
	}
	return &EtcdRegistry{Client: cli, config: c}, nil
}

func (e *EtcdRegistry) Name() string {
	return "etcd"
}

//
//func (e *EtcdRegistry) Init(ctx context.Context) (registry.Registryer, error) {
//	client, err := NewEtcdRegistry()
//	return client, err
//}

func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) error {
	panic("implement me")
}

func (e *EtcdRegistry) UnRegister(ctx context.Context, service *registry.Service) error {
	panic("implement me")
}

func (e *EtcdRegistry) Get(ctx context.Context, service *registry.Service) error {
	panic("implement me")
}
