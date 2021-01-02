package consul

import "errors"

var ConsulNotFoundErr = errors.New("The client could not be found")
var ConsulRegistionErr = errors.New("Registion failed")

var NacosWeightErr = errors.New("Weight must be lager than 0")
var NacosParamErr = errors.New("Registion failed")
var NacosDeregistionErr = errors.New("Deregistion failed")
var NacosGetServiceErr = errors.New("Nacos get services failed")
var NacosSubscribeErr = errors.New("Nacos subscription failed")
var NacosUnsubscribeErr = errors.New("Nacos Unsubscription failed")

//DeregisterInstance
