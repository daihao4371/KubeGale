package initialize

import (
	"KubeGale/router"
	"KubeGale/router/cloudCmdb"

	"github.com/gin-gonic/gin"
)

func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	{
		// 注册IM通信路由
		imRouter := router.RouterGroupApp.Im
		imRouter.InitNotificationRouter(privateGroup)
	}

	holder(publicGroup, privateGroup) // 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
	{
		cmdbRouter := router.RouterGroupApp.Cmdb
		cmdbRouter.InitCmdbProjectsRouter(privateGroup, publicGroup)
		cmdbRouter.InitCmdbHostsRouter(privateGroup, publicGroup)
		cmdbRouter.InitBatchOperationsRouter(privateGroup, publicGroup)
	}

	{
		// 注册云资源路由
		cloudCmdbRouter := cloudCmdb.RouterGroupApp
		cloudCmdbRouter.InitCloudPlatformRouter(privateGroup)  // 云平台路由
		cloudCmdbRouter.InitCloudRegionRouter(privateGroup)    // 云区域路由
		cloudCmdbRouter.InitVirtualMachineRouter(privateGroup) // 云主机路由
		cloudCmdbRouter.InitLoadBalancerRouter(privateGroup)   // 负载均衡路由
		cloudCmdbRouter.InitRDSRouter(privateGroup)            // RDS路由
	}
}
