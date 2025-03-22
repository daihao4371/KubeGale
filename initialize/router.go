package initialize

import (
	"KubeGale/router"
	"github.com/gin-gonic/gin"
	"time"
)

func Routers() *gin.Engine {
	Router := gin.New()
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
		// 初始化用户相关路由
		systemRouter.InitUserRouter(PublicGroup)
		systemRouter.InitUserRouter(PrivateGroup)
		// 用户路由包含了登录等公共接口
	}

	// 添加健康检查路由
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	return Router
}
