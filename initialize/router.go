package initialize

import (
	"KubeGale/router"
	"github.com/gin-gonic/gin"
	"time"
)

// 这里只展示需要修改的部分
func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	// 初始化公共路由组
	PublicGroup := Router.Group("/api")

	// 初始化需要认证的路由组
	PrivateGroup := Router.Group("/api")

	// 注册系统路由
	systemRouter := router.RouterGroupApp.System
	{
		// 初始化用户相关路由 - 只在一个路由组中注册
		systemRouter.InitUserRouter(PublicGroup)
		systemRouter.InitMenuRouter(PrivateGroup)               // 注册menu路由
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)   // 注册功能api路由
		systemRouter.InitSysOperationRecordRouter(PrivateGroup) // 操作记录
		systemRouter.InitAuthorityRouter(PrivateGroup)          // 权限分配
		systemRouter.InitRoleRouter(PrivateGroup)               // 权限分配

		// 不要在 PrivateGroup 中重复注册相同的路由
		// systemRouter.InitUserRouter(PrivateGroup)
	}

	// 用户路由包含了登录等公共接口
	// 添加健康检查路由
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	return Router
}
