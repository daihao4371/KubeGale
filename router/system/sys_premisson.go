package system

import (
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("permissions").Use(middleware.OperationRecord())
	//authorityRouterWithoutRecord := Router.Group("permissions")
	// 需要记录操作日志的路由
	{
		authorityRouter.POST("/user/assign", authorityApi.AssignUserRole)   // 为单个用户分配角色和权限
		authorityRouter.POST("/users/assign", authorityApi.AssignUsersRole) // 批量为用户分配角色和权限
	}
}
