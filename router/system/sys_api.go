package system

import (
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	// API相关路由组
	apiRouter := Router.Group("apis").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := RouterPub.Group("apis")

	// API管理相关路由
	{
		// 需要记录操作日志的API路由
		apiRouter.POST("/create", apiRouterApi.CreateAPI) // 创建API
		apiRouter.POST("/update", apiRouterApi.UpdateAPI) // 更新API
		apiRouter.DELETE("/:id", apiRouterApi.DeleteAPI)  // 删除API

		// 不需要记录操作日志的API路由
		apiRouterWithoutRecord.GET("/list", apiRouterApi.ListApis) // API列表
	}
}
