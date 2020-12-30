package registry

import (
	"context"
	"fmt"
	"sync"

	vo "github.com/duckbb/registry/base"
)

type RegistryType string

const (
	NACOS     RegistryType = "nacos"
	ZOOKEEPER RegistryType = "zookeeper"
)

type pluninRegistry struct {
	sync.Mutex
	plunins             map[RegistryType]vo.Registryer
	currentRegisterType RegistryType
}

var PR *pluninRegistry
var once sync.Once

func GetRegistryMgr() *pluninRegistry {
	once.Do(func() {
		PR = &pluninRegistry{plunins: make(map[RegistryType]vo.Registryer)}
	})
	return PR
}

//set current RegisterType
func (p *pluninRegistry) SetRegisterType(name RegistryType) {
	p.Lock()
	defer p.Unlock()
	p.currentRegisterType = name
}

//plunin init
func (p *pluninRegistry) InitPluninRegistry(ctx context.Context, name RegistryType, r vo.Registryer) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.plunins[name]; ok {
		return fmt.Errorf("%s plun-in has registered", name)
	}
	p.plunins[name] = r
	return nil
}

//registry exist
func (p *pluninRegistry) GetRegister(name ...RegistryType) (vo.Registryer, error) {
	rType := p.currentRegisterType
	if len(name) > 0 {
		rType = name[0]
	}
	r, ok := p.plunins[rType]
	if !ok {
		return nil, fmt.Errorf("%s Plug-in not registered", name)
	}
	return r, nil
}

//register base
func (p *pluninRegistry) Register(ctx context.Context, service *vo.Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.GetRegister(name...)
	if err != nil {
		return err
	}
	return r.Register(ctx, service)
}

//unregister base
func (p *pluninRegistry) UnRegister(ctx context.Context, service *vo.Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.GetRegister(name...)
	if err != nil {
		return err
	}
	return r.UnRegister(ctx, service)
}

//get services
func (p *pluninRegistry) Get(ctx context.Context, service *vo.Service, name ...RegistryType) ([]*vo.Service, error) {
	p.Lock()
	defer p.Unlock()
	r, err := p.GetRegister(name...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, service)
}

func (p *pluninRegistry) SubscribeService(ctx context.Context, service *vo.Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.GetRegister(name...)
	if err != nil {
		return err
	}
	return r.SubscribeService(ctx, service)
}
