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
	{
		// 注册IM通信路由
		imRouter := router.RouterGroupApp.Im
		imRouter.InitNotificationRouter(privateGroup)
	}
}
