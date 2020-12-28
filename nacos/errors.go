package nacos

import "errors"

var NacosNotFoundErr = errors.New("The client could not be found")
var NacosWeightErr = errors.New("Weight must be lager than 0")
var NacosRegistionErr = errors.New("registion failed")
var NacosParamErr = errors.New("registion failed")
var NacosDeregistionErr = errors.New("deregistion failed")

//DeregisterInstance
