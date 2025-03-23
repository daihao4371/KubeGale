package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct{}

// AssignUserRole 为单个用户分配角色和权限
func (h *AuthorityApi) AssignUserRole(c *gin.Context) {
	var req system.AssignUserRoleRequest
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 参数验证
	if req.UserId <= 0 {
		response.FailWithMessage("用户ID无效", c)
		return
	}

	// 调用服务层分配角色和权限
	err := authorityService.AssignRoleToUser(req.UserId, req.RoleIds, req.ApiIds)
	if err != nil {
		global.KUBEGALE_LOG.Error("分配角色和权限失败", zap.Error(err), zap.Int("userId", req.UserId))
		response.FailWithMessage("分配角色和权限失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithDetailed(gin.H{
		"userId":  req.UserId,
		"roleIds": req.RoleIds,
		"apiIds":  req.ApiIds,
	}, "分配角色和权限成功", c)
}

// AssignUsersRole 批量为用户分配角色和权限
func (h *AuthorityApi) AssignUsersRole(c *gin.Context) {
	var req system.AssignUsersRoleRequest
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 调用服务层批量分配角色和权限
	err := authorityService.AssignRoleToUsers(req.UserIds, req.RoleIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		global.KUBEGALE_LOG.Error("批量分配角色和权限失败")
		return
	}
	response.OkWithMessage("批量分配角色和权限成功", c)
}
