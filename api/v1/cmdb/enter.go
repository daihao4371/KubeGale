package cmdb

import "KubeGale/service"

type ApiGroup struct {
	CmdbProjectsApi
	CmdbHostsApi
	BatchOperationsApi
}

var cmdbProjectsService = service.ServiceGroupApp.CmdbServiceGroup.CmdbProjectsService
var cmdbHostsService = service.ServiceGroupApp.CmdbServiceGroup.CmdbHostsService
var batchOperationsService = service.ServiceGroupApp.CmdbServiceGroup.BatchOperationsService

