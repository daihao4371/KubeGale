package system

import (
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type RoleRouter struct{}

func (s *RoleRouter) InitRoleRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("roles").Use(middleware.OperationRecord())
	//authorityRouterWithoutRecord := Router.Group("permissions")
	// 需要记录操作日志的路由
	{
		authorityRouter.POST("/list", roleApi.ListRoles)       // 为单个用户分配角色和权限
		authorityRouter.POST("/create", roleApi.CreateRole)    // 批量为用户分配角色和权限
		authorityRouter.POST("/update", roleApi.UpdateRole)    // 批量为用户分配角色和权限
		authorityRouter.DELETE("/:id", roleApi.DeleteRole)     // 批量为用户分配角色和权限
		authorityRouter.GET("/user/:id", roleApi.GetUserRoles) // 批量为用户分配角色和权限
		authorityRouter.GET("/:id", roleApi.GetRoles)          // 批量为用户分配角色和权限
	}
}
