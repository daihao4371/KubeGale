package initialize

import (
	"KubeGale/middleware"
	"KubeGale/router"
	"github.com/gin-gonic/gin"
	"time"
)

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	systemRouter := router.RouterGroupApp.System

	// 添加SQL日志中间件
	Router.Use(middleware.SQLLogMiddleware())

	// 初始化公共路由组
	PublicGroup := Router.Group("/api")

	// 初始化需要认证的路由组
	PrivateGroup := Router.Group("/api")

	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		//systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}

	// 注册系统路由
	{
		// 初始化用户相关路由 - 只在一个路由组中注册
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)   // 注册功能api路由
		systemRouter.InitJwtRouter(PrivateGroup)                // jwt相关路由
		systemRouter.InitUserRouter(PrivateGroup)               // 注册用户路由
		systemRouter.InitMenuRouter(PrivateGroup)               // 注册menu路由
		systemRouter.InitCasbinRouter(PrivateGroup)             // 权限相关路由
		systemRouter.InitAuthorityRouter(PrivateGroup)          // 注册角色路由
		systemRouter.InitSysOperationRecordRouter(PrivateGroup) // 操作记录
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup) // 按钮权限管理
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
