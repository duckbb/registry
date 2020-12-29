package nacos

import "errors"

var NacosNotFoundErr = errors.New("The client could not be found")
var NacosWeightErr = errors.New("Weight must be lager than 0")
var NacosRegistionErr = errors.New("Registion failed")
var NacosParamErr = errors.New("Registion failed")
var NacosDeregistionErr = errors.New("Deregistion failed")
var NacosGetServiceErr = errors.New("Nacos get services failed")

//DeregisterInstance
