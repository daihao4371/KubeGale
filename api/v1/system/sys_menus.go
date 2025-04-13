package system

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	systemReq "KubeGale/model/system/request"
	systemRes "KubeGale/model/system/response"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityMenuApi struct{}

// AuthorityMenu 获取用户动态路由
func (a *AuthorityMenuApi) GetMenu(c *gin.Context) {
	menus, err := menuService.GetMenuTree(utils.GetUserAuthorityId(c))
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	if menus == nil {
		menus = []system.SysMenu{}
	}
	response.OkWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "获取成功", c)
}

// GetBaseMenuTree 获取用户动态路由
func (a *AuthorityMenuApi) GetBaseMenuTree(c *gin.Context) {
	authority := utils.GetUserAuthorityId(c)
	menus, err := menuService.GetBaseMenuTree(authority)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenusResponse{Menus: menus}, "获取成功", c)
}

// AddMenuAuthority 增加menu和角色关联关系
func (a *AuthorityMenuApi) AddMenuAuthority(c *gin.Context) {
	var authorityMenu systemReq.AddMenuAuthorityInfo
	err := c.ShouldBindJSON(&authorityMenu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(authorityMenu, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminAuthorityID := utils.GetUserAuthorityId(c)
	if err := menuService.AddMenuAuthority(authorityMenu.Menus, adminAuthorityID, authorityMenu.AuthorityId); err != nil {
		global.KUBEGALE_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// GetMenuAuthority 获取指定角色menu
func (a *AuthorityMenuApi) GetMenuAuthority(c *gin.Context) {
	var param request.GetAuthorityId
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menus, err := menuService.GetMenuAuthority(&param)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"menus": menus}, "获取成功", c)
}

// AddBaseMenu 新增菜单
func (a *AuthorityMenuApi) AddBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = menuService.AddBaseMenu(menu)
	if err != nil {
		global.KUBEGALE_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// DeleteBaseMenu 删除菜单
func (a *AuthorityMenuApi) DeleteBaseMenu(c *gin.Context) {
	var menu request.GetById
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.DeleteBaseMenu(menu.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateBaseMenu  更新菜单
func (a *AuthorityMenuApi) UpdateBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.UpdateBaseMenu(menu)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetBaseMenuById 根据id获取菜单
func (a *AuthorityMenuApi) GetBaseMenuById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menu, err := baseMenuService.GetBaseMenuById(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenuResponse{Menu: menu}, "获取成功", c)
}

// GetMenuList 分页获取基础menu列表
func (a *AuthorityMenuApi) GetMenuList(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	menuList, err := menuService.GetInfoList(authorityID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(menuList, "获取成功", c)
}
