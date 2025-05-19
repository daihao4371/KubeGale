package cloudCmdb

import (
	api "KubeGale/api/v1/cloudCmdb"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type CloudVirtualMachineRouter struct {
}

func (c *CloudVirtualMachineRouter) InitVirtualMachineRouter(Router *gin.RouterGroup) {
	cloudVirtualMachineRouter := Router.Group("virtualMachine").Use(middleware.OperationRecord())
	cloudVirtualMachineRouterWithoutRecord := Router.Group("virtualMachine")
	cloudVirtualMachineApi := api.ApiGroupApp.CloudVirtualMachineApi
	{
		cloudVirtualMachineRouter.POST("sync", cloudVirtualMachineApi.CloudVirtualMachineSync) // 同步云主机
		cloudVirtualMachineRouter.POST("tree", cloudVirtualMachineApi.CloudVirtualMachineTree) // 目录树
	}

	{
		cloudVirtualMachineRouterWithoutRecord.POST("list", cloudVirtualMachineApi.CloudVirtualMachineList) // 分页获取列表
	}
}
