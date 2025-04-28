package service

import (
	"KubeGale/service/im"
	"KubeGale/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
	ImServiceGroup     im.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)