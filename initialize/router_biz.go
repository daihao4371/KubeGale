package initialize

import (
	"KubeGale/router"

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

}
