package initialize

import (
	"KubeGale/global"
	"KubeGale/middleware"
	"KubeGale/router"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	systemRouter := router.RouterGroupApp.System
	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面3行注释
	// Router.Static("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	// Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	// global.KUBEGALE_LOG.Info("use middleware cors")
	//docs.SwaggerInfo.BasePath = global.KUBEGALE_CONFIG.System.RouterPrefix
	//Router.GET(global.KUBEGALE_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.KUBEGALE_LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group(global.KUBEGALE_CONFIG.System.RouterPrefix)
	PrivateGroup := Router.Group(global.KUBEGALE_CONFIG.System.RouterPrefix)

	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		//systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}

	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)   // 注册功能api路由
		systemRouter.InitJwtRouter(PrivateGroup)                // jwt相关路由
		systemRouter.InitUserRouter(PrivateGroup)               // 注册用户路由
		systemRouter.InitMenuRouter(PrivateGroup)               // 注册menu路由
		systemRouter.InitCasbinRouter(PrivateGroup)             // 权限相关路由
		systemRouter.InitAuthorityRouter(PrivateGroup)          // 注册角色路由
		systemRouter.InitSysOperationRecordRouter(PrivateGroup) // 操作记录
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup) // 按钮权限管理

	}

	//插件路由安装
	//InstallPlugin(PrivateGroup, PublicGroup, Router)

	// 注册业务路由
	initBizRouter(PrivateGroup, PublicGroup)

	global.KUBEGALE_ROUTERS = Router.Routes()

	global.KUBEGALE_LOG.Info("router register success")
	return Router
}
