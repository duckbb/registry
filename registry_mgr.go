package registry

import (
	"context"
	"fmt"
	"sync"
)

type RegistryType string

const (
	NACOS     RegistryType = "nacos"
	ZOOKEEPER RegistryType = "zookeeper"
)

type pluninRegistry struct {
	sync.Mutex
	plunins             map[RegistryType]Registryer
	currentRegisterType RegistryType
}

var PR *pluninRegistry
var once sync.Once

func GetRegistryMgr() *pluninRegistry {
	once.Do(func() {
		PR = &pluninRegistry{plunins: make(map[RegistryType]Registryer)}
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
func (p *pluninRegistry) InitPluninRegistry(ctx context.Context, name RegistryType, r Registryer) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.plunins[name]; ok {
		return fmt.Errorf("%s plun-in has registered", name)
	}
	p.plunins[name] = r
	return nil
}

//registry exist
func (p *pluninRegistry) getRegister(service *Service, name ...RegistryType) (Registryer, error) {
	var rType RegistryType = p.currentRegisterType
	if len(name) > 0 {
		rType = name[0]
	}
	r, ok := p.plunins[rType]
	if ok {
		return nil, fmt.Errorf("%s plun-in has registered", name)
	}
	return r, nil
}

//register service
func (p *pluninRegistry) Register(ctx context.Context, service *Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.getRegister(service, name...)
	if err != nil {
		return err
	}
	return r.Register(ctx, service)
}

//unregister service
func (p *pluninRegistry) UnRegister(ctx context.Context, service *Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.getRegister(service, name...)
	if err != nil {
		return err
	}
	return r.UnRegister(ctx, service)
}

//get services
func (p *pluninRegistry) Get(ctx context.Context, service *Service, name ...RegistryType) ([]*Service, error) {
	p.Lock()
	defer p.Unlock()
	r, err := p.getRegister(service, name...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, service)
}

func (p *pluninRegistry) SubscribeService(ctx context.Context, service *Service, name ...RegistryType) error {
	p.Lock()
	defer p.Unlock()
	r, err := p.getRegister(service, name...)
	if err != nil {
		return err
	}
	return r.SubscribeService(ctx, service)
}
