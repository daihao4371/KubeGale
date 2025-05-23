package cloudCmdb

type ServiceGroup struct {
	CloudPlatformService
	CloudRegionService
	CloudVirtualMachineService
	CloudLoadBalancerService
	CloudRDSService
	CloudStatsService // Add this line
}

var ServiceGroupApp = new(ServiceGroup)
