package system

import (
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (s *MenuRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	// 创建路由组
	menuRouter := Router.Group("menus").Use(middleware.OperationRecord())
	menuRouterWithoutRecord := Router.Group("menus")

	// 需要记录操作日志的路由
	{
		menuRouter.POST("/create", authorityMenuApi.CreateMenu)             // 创建菜单
		menuRouter.POST("/update", authorityMenuApi.UpdateMenu)             // 更新菜单
		menuRouter.DELETE("/:id", authorityMenuApi.DeleteMenu)              // 删除菜单
		menuRouter.POST("/update_related", authorityMenuApi.UpdateUserMenu) // 更新用户菜单关联
	}

	// 不需要记录操作日志的路由
	{
		menuRouterWithoutRecord.POST("/list", authorityMenuApi.ListMenus) // 获取菜单列表
	}

	return menuRouter
}
