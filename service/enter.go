package service

import (
	"KubeGale/service/cmdb"
	"KubeGale/service/im"
	"KubeGale/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
	ImServiceGroup     im.ServiceGroup
	CmdbServiceGroup   cmdb.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
