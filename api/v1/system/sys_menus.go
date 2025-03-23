package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type AuthorityMenuApi struct{}

func (m *AuthorityMenuApi) ListMenus(c *gin.Context) {
	var req system.ListMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}
	// 调用service层获取菜单列表
	menus, total, err := menuService.GetMenus(req.PageNumber, req.PageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取菜单列表失败", zap.Error(err))
		response.FailWithMessage("获取菜单列表失败: "+err.Error(), c)
		return
	}

	// 返回菜单列表和总数
	response.OkWithDetailed(gin.H{
		"list":     menus,
		"total":    total,
		"page":     req.PageNumber,
		"pageSize": req.PageSize,
	}, "获取菜单列表成功", c)
}

func (m *AuthorityMenuApi) CreateMenu(c *gin.Context) {
	var req system.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menu := &system.Menu{
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		ParentID:  req.ParentId,
		Hidden:    req.Hidden,
		RouteName: req.RouteName,
		Redirect:  req.Redirect,
		Meta:      req.Meta,
		Children:  req.Children,
	}

	if err := menuService.CreateMenu(menu); err != nil {
		response.FailWithMessage(err.Error(), c)
		global.KUBEGALE_LOG.Error("创建菜单失败")
		return
	}
	response.OkWithMessage("创建菜单成功", c)
}

// DeleteMenu 删除菜单
func (m *AuthorityMenuApi) DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := menuService.DeleteMenu(id); err != nil {
		response.FailWithMessage("删除菜单失败", c)
		global.KUBEGALE_LOG.Error("删除菜单失败")
		return
	}

	response.OkWithMessage("删除菜单成功", c)
}

// UpdateMenu 更新菜单
func (m *AuthorityMenuApi) UpdateMenu(c *gin.Context) {
	var req system.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	menu := &system.Menu{
		ID:        req.Id,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		ParentID:  req.ParentId,
		Hidden:    req.Hidden,
		RouteName: req.RouteName,
	}

	if err := menuService.UpdateMenu(menu); err != nil {
		response.FailWithMessage(err.Error(), c)
		global.KUBEGALE_LOG.Error("更新菜单失败")
		return
	}
	response.OkWithMessage("更新菜单失败", c)
}

// AddUserMenu 添加用户菜单关联
func (m *AuthorityMenuApi) UpdateUserMenu(c *gin.Context) {
	var req system.UpdateUserMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := menuService.UpdateUserMenu(req.UserId, req.MenuIds); err != nil {
		response.FailWithMessage("更新用户菜单关联失败", c)
		global.KUBEGALE_LOG.Error("更新用户菜单关联失败")
		return
	}
	response.OkWithMessage("更新用户菜单关联成功", c)
}
