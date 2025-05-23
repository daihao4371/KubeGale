package service

import (
	"KubeGale/service/cmdb"
	"KubeGale/service/cloudCmdb" // Added this import
	"KubeGale/service/im"
	"KubeGale/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup    system.ServiceGroup
	ImServiceGroup        im.ServiceGroup
	CmdbServiceGroup      cmdb.ServiceGroup
	CloudCmdbServiceGroup cloudCmdb.ServiceGroup // Added this line
}

var ServiceGroupApp = new(ServiceGroup)
