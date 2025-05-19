package cloudCmdb

type ServiceGroup struct {
	CloudPlatformService
	CloudRegionService
	CloudVirtualMachineService
	CloudLoadBalancerService
	CloudRDSService
}

var ServiceGroupApp = new(ServiceGroup)
