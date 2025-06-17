package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthorityApiApi 角色-API权限相关接口
// 建议注册到 ApiGroup 结构体

type AuthorityApiApi struct{}

var AuthorityApiApiApp = new(AuthorityApiApi)

// SetApisForAuthority 为角色分配API权限
func (a *AuthorityApiApi) SetApisForAuthority(c *gin.Context) {
	var req struct {
		AuthorityId uint   `json:"authorityId" binding:"required"`
		ApiIds      []uint `json:"apiIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		global.KUBEGALE_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := authorityApiService.SetApisForAuthority(req.AuthorityId, req.ApiIds); err != nil {
		global.KUBEGALE_LOG.Error("分配API权限失败!", zap.Error(err))
		response.FailWithMessage("分配API权限失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("分配API权限成功", c)
}

// GetApisByAuthority 查询角色拥有的API权限
func (a *AuthorityApiApi) GetApisByAuthority(c *gin.Context) {
	var req struct {
		AuthorityId uint `json:"authorityId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		global.KUBEGALE_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	apis, err := authorityApiService.GetApisByAuthority(req.AuthorityId)
	if err != nil {
		global.KUBEGALE_LOG.Error("查询API权限失败!", zap.Error(err))
		response.FailWithMessage("查询API权限失败: "+err.Error(), c)
		return
	}
	response.OkWithData(apis, c)
}
