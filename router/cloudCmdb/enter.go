package cloudCmdb

type RouterGroup struct {
	CloudPlatformRouter
	CloudRegionRouter
	CloudVirtualMachineRouter
	CloudLoadBalancerRouter
	CloudRDSRouter
	CloudStatsRouter // Add this line
}

var RouterGroupApp = new(RouterGroup)
