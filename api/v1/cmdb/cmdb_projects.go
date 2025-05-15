package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	"KubeGale/model/common/response"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CmdbProjectsApi struct{}

// CreateCmdbProjects 创建cmdbProjects表
func (cmdbProjectsApi *CmdbProjectsApi) CreateCmdbProjects(c *gin.Context) {
	var cmdbProjects cmdb.CmdbProjects
	err := c.ShouldBindJSON(&cmdbProjects)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbProjects.CreatedBy = utils.GetUserID(c)
	err = cmdbProjectsService.CreateCmdbProjects(&cmdbProjects)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCmdbProjects 删除cmdbProjects表
func (cmdbProjectsApi *CmdbProjectsApi) DeleteCmdbProjects(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := cmdbProjectsService.DeleteCmdbProjects(ID, userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCmdbProjectsByIds 批量删除cmdbProjects表
func (cmdbProjectsApi *CmdbProjectsApi) DeleteCmdbProjectsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := cmdbProjectsService.DeleteCmdbProjectsByIds(IDs, userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCmdbProjects 更新cmdbProjects表
func (cmdbProjectsApi *CmdbProjectsApi) UpdateCmdbProjects(c *gin.Context) {
	var cmdbProjects cmdb.CmdbProjects
	err := c.ShouldBindJSON(&cmdbProjects)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbProjects.UpdatedBy = utils.GetUserID(c)
	err = cmdbProjectsService.UpdateCmdbProjects(cmdbProjects)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCmdbProjects 用id查询cmdbProjects表
func (cmdbProjectsApi *CmdbProjectsApi) FindCmdbProjects(c *gin.Context) {
	ID := c.Query("ID")
	recmdbProjects, err := cmdbProjectsService.GetCmdbProjects(ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(recmdbProjects, c)
}

// GetCmdbProjectsList 分页获取cmdbProjects表列表
func (cmdbProjectsApi *CmdbProjectsApi) GetCmdbProjectsList(c *gin.Context) {
	var pageInfo cmdbReq.CmdbProjectsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := cmdbProjectsService.GetCmdbProjectsInfoList(pageInfo)
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
