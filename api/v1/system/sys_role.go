package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleApi struct{}

// ListRoles 获取角色列表
func (r *RoleApi) ListRoles(c *gin.Context) {
	var req system.ListRolesRequest

	// 尝试绑定JSON，如果失败则使用默认值
	if err := c.ShouldBindJSON(&req); err != nil {
		global.KUBEGALE_LOG.Warn("请求参数绑定失败，将使用默认分页参数", zap.Error(err))
		// 使用默认值
		req.PageNumber = 1
		req.PageSize = 10
	}

	// 确保分页参数有效
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用service获取角色列表
	roles, total, err := roleService.ListRoles(req.PageNumber, req.PageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取角色列表失败", zap.Error(err))
		response.FailWithMessage("获取角色列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  roles,
		"total": total,
	}, "获取角色列表成功", c)
}

// CreateRole 创建角色
func (r *RoleApi) CreateRole(c *gin.Context) {
	var req system.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 构建角色对象
	role := &system.Role{
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsDefault:   req.IsDefault,
	}

	// 创建角色并分配权限
	err := roleService.CreateRole(role, req.ApiIds)
	if err != nil {
		response.FailWithMessage("创建角色并分配权限失败", c)
		global.KUBEGALE_LOG.Error("创建角色并分配权限失败")
		return
	}

	response.OkWithMessage("创建角色并分配权限成功", c)
}

// UpdateRole 更新角色
func (r *RoleApi) UpdateRole(c *gin.Context) {
	var req system.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 构建角色对象
	role := &system.Role{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsDefault:   req.IsDefault,
	}

	// 更新角色基本信息和权限（在service层已经处理了权限分配）
	if err := roleService.UpdateRole(role, req.ApiIds); err != nil {
		global.KUBEGALE_LOG.Error("更新角色信息失败", zap.Error(err), zap.Int("roleId", req.Id))
		response.FailWithMessage("更新角色信息失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithDetailed(gin.H{
		"id":     req.Id,
		"name":   req.Name,
		"apiIds": req.ApiIds,
	}, "更新角色成功", c)
}

// DeleteRole 删除角色
func (r *RoleApi) DeleteRole(c *gin.Context) {
	// 从URL参数中获取角色ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = roleService.DeleteRole(id)
	if err != nil {
		response.FailWithMessage("删除角色失败", c)
		global.KUBEGALE_LOG.Error("删除角色失败")
		return
	}
	response.OkWithMessage("删除角色成功", c)
}

// UpdateUserRole 更新用户角色
func (r *RoleApi) UpdateUserRole(c *gin.Context) {
	var req system.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 分配用户角色和权限
	err := authorityService.AssignRoleToUser(req.UserId, req.RoleIds, req.ApiIds)
	if err != nil {
		response.FailWithMessage("分配用户角色和权限失败", c)
		global.KUBEGALE_LOG.Error("分配用户角色和权限失败")
		return
	}

	response.OkWithMessage("分配用户角色和权限成功", c)
}

// GetUserRoles 获取用户角色
func (r *RoleApi) GetUserRoles(c *gin.Context) {
	// 从URL参数中获取用户ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	role, err := roleService.GetUserRole(id)
	if err != nil {
		response.FailWithMessage("获取用户角色失败", c)
		global.KUBEGALE_LOG.Error("获取用户角色失败")
		return
	}
	response.OkWithDetailed(role, "获取用户角色失败", c)
}

// GetRoles 获取角色详情
func (r *RoleApi) GetRoles(c *gin.Context) {
	// 从URL参数中获取角色ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	role, err := roleService.GetRole(id)
	if err != nil {
		response.FailWithMessage("获取角色详情失败", c)
		global.KUBEGALE_LOG.Error("获取角色详情失败")
		return
	}
	response.OkWithDetailed(role, "获取角色详情成功", c)
}
