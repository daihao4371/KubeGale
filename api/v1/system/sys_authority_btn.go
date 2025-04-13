package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system/request"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type AuthorityBtnApi struct{}

// GetAuthorityBtn 获取权限按钮
func (a *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(res, "查询成功", c)
}

// SetAuthorityBtn 设置权限按钮
func (a *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityBtnService.SetAuthorityBtn(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("分配失败!", zap.Error(err))
		response.FailWithMessage("分配失败", c)
		return
	}
	response.OkWithMessage("分配成功", c)
}

// CanRemoveAuthorityBtn 设置权限按钮
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	err := authorityBtnService.CanRemoveAuthorityBtn(id)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
