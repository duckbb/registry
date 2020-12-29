package registry

//
//import (
//	"context"
//	"fmt"
//	"sync"
//)
//
//type pluninRegistry struct {
//	sync.Mutex
//	plunins map[string]Registryer
//}
//
//var PR *pluninRegistry
//var once sync.Once
//
//func GetRegistryMgr() *pluninRegistry {
//	once.Do(func() {
//		PR = &pluninRegistry{plunins: make(map[string]Registryer)}
//	})
//	return PR
//}
//
////plunin regisrer
//func (p *pluninRegistry) RegisterPlunin(r Registryer) error {
//	p.Lock()
//	defer p.Unlock()
//	if _, ok := p.plunins[r.Name()]; ok {
//		return fmt.Errorf("%s had register", r.Name())
//	}
//	p.plunins[r.Name()] = r
//	return nil
//}
//
////plunin init
//func (p *pluninRegistry) InitPluninRegistry(ctx context.Context, name string, opts ...Option) error {
//	p.Lock()
//	defer p.Unlock()
//	if _, ok := p.plunins[name]; ok {
//		return fmt.Errorf("%s plun-in not register", name)
//	}
//	err := p.plunins[name].Init(ctx, opts...)
//	if err != nil {
//		return fmt.Errorf("plun-in init fail,err:%s", err)
//	}
//	return nil
//}
