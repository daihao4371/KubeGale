package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	"KubeGale/model/common/response"
	"KubeGale/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CmdbHostsApi struct{}

// CreateCmdbHosts 创建cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) CreateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.CreateCmdbHosts(&cmdbHosts)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// AuthenticationCmdbHosts 验证cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) AuthenticationCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.SSHTestCmdbHosts(&cmdbHosts)
	if err != nil {
		if err.Error() == "auth failed" {
			response.Result(177, nil, "auth failed", c)
			return
		}
		global.KUBEGALE_LOG.Error("验证失败!", zap.Error(err))
		response.FailWithMessage("验证失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("验证成功", c)
}

// DeleteCmdbHosts 删除cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHosts(c *gin.Context) {
	var req cmdbReq.DeleteCmdbHostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHosts(fmt.Sprintf("%d", req.ID), userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCmdbHostsByIds 批量删除cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHostsByIds(c *gin.Context) {
	var req cmdbReq.DeleteCmdbHostsIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 将 uint 数组转换为 string 数组
	ids := make([]string, len(req.IDs))
	for i, id := range req.IDs {
		ids[i] = fmt.Sprintf("%d", id)
	}

	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHostsByIds(ids, userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCmdbHosts 更新cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) UpdateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.UpdatedBy = utils.GetUserID(c)
	err = cmdbHostsService.UpdateCmdbHosts(cmdbHosts)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCmdbHosts 用id查询cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) FindCmdbHosts(c *gin.Context) {
	ID := c.Query("id")
	recmdbHosts, err := cmdbHostsService.GetCmdbHosts(ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(recmdbHosts, c)
}

// GetCmdbHostsList 获取cmdbHosts表列表
func (cmdbHostsApi *CmdbHostsApi) GetCmdbHostsList(c *gin.Context) {
	var pageInfo cmdbReq.CmdbHostsSearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		global.KUBEGALE_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	list, total, err := cmdbHostsService.GetCmdbHostsInfoList(pageInfo)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ImportHosts 根据模板批量创建主机
func (cmdbHostsApi *CmdbHostsApi) ImportHosts(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file") // "file" 是前端上传文件的字段名
	if err != nil {
		response.FailWithMessage("获取文件失败: "+err.Error(), c)
		return
	}
	projectIdStr := c.PostForm("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		response.FailWithMessage("无效的 projectId: "+err.Error(), c)
		return
	}

	// 保存上传的文件到临时目录
	dst := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.FailWithMessage("保存文件失败: "+err.Error(), c)
		return
	}

	// 调用服务层处理上传逻辑
	if err := cmdbHostsService.ImportHosts(dst, projectId); err != nil {
		global.KUBEGALE_LOG.Error("导入失败!", zap.Error(err))
		response.FailWithMessage("导入失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("导入成功", c)
}
