package cloudCmdb

import (
	service "KubeGale/service/cloudCmdb"
)

type ApiGroup struct {
	CloudPlatformApi
	CloudRegionApi
	CloudVirtualMachineApi
	CloudLoadBalancerApi
	CloudRDSApi
	CloudStatsApi // Add this line
}

var ApiGroupApp = new(ApiGroup)

var (
	cloudPlatformService       = service.ServiceGroupApp.CloudPlatformService
	cloudVirtualMachineService = service.ServiceGroupApp.CloudVirtualMachineService
	cloudRegionService         = service.ServiceGroupApp.CloudRegionService
	cloudLoadBalancerService   = service.ServiceGroupApp.CloudLoadBalancerService
	cloudRDSService            = service.ServiceGroupApp.CloudRDSService
)
