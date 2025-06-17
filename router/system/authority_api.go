package system

import (
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type AuthorityApiRouter struct{}

// InitAuthorityApiRouter 注册角色-API权限相关路由
func (s *AuthorityApiRouter) InitAuthorityApiRouter(Router *gin.RouterGroup) {
	authorityApiRouter := Router.Group("authorityApi").Use(middleware.OperationRecord())
	{
		authorityApiRouter.POST("/setApisForAuthority", authorityApiApi.SetApisForAuthority) // 为角色分配API权限
		authorityApiRouter.POST("/getApisByAuthority", authorityApiApi.GetApisByAuthority)   // 查询角色API权限
	}
}
