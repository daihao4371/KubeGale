package cmdb

import (
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type CmdbHostsRouter struct{}

// InitCmdbHostsRouter 初始化 cmdbHosts表 路由信息
func (s *CmdbHostsRouter) InitCmdbHostsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	cmdbHostsRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	cmdbHostsRouterWithoutRecord := Router.Group("cmdb")
	{
		cmdbHostsRouter.POST("hosts", cmdbHostsApi.CreateCmdbHosts)                        // 新建cmdbHosts表
		cmdbHostsRouter.DELETE("hosts", cmdbHostsApi.DeleteCmdbHosts)                      // 删除cmdbHosts表
		cmdbHostsRouter.DELETE("hostsByIds", cmdbHostsApi.DeleteCmdbHostsByIds)            // 批量删除cmdbHosts表
		cmdbHostsRouter.PUT("hosts", cmdbHostsApi.UpdateCmdbHosts)                         // 更新cmdbHosts表
		cmdbHostsRouter.POST("hosts/authentication", cmdbHostsApi.AuthenticationCmdbHosts) // ssh认证主机
		cmdbHostsRouter.POST("hosts/import", cmdbHostsApi.ImportHosts)                     // 批量导入
	}
	{
		cmdbHostsRouterWithoutRecord.GET("hostsById", cmdbHostsApi.FindCmdbHosts)     // 根据ID获取cmdbHosts表
		cmdbHostsRouterWithoutRecord.POST("hostsList", cmdbHostsApi.GetCmdbHostsList) // 获取cmdbHosts表列表
	}

	// 将 WebSocket 终端路由移到公共路由组
	cmdbHostsPublicRouter := PublicRouter.Group("cmdb")
	{
		cmdbHostsPublicRouter.GET("hosts/terminal", cmdbHostsApi.WebTerminal) // 处理Web终端连接
	}
}
