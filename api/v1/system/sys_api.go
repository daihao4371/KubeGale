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

type SystemApiApi struct{}

// CreateApi
// @Summary   创建基础api
func (s *SystemApiApi) CreateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.CreateApi(api)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// SyncApi
// @Summary   同步API
func (s *SystemApiApi) SyncApi(c *gin.Context) {
	newApis, deleteApis, ignoreApis, err := apiService.SyncApi()
	if err != nil {
		global.KUBEGALE_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	response.OkWithData(gin.H{
		"newApis":    newApis,
		"deleteApis": deleteApis,
		"ignoreApis": ignoreApis,
	}, c)
}

// GetApiGroups
// @Summary   获取API分组
func (s *SystemApiApi) GetApiGroups(c *gin.Context) {
	groups, apiGroupMap, err := apiService.GetApiGroups()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(gin.H{
		"groups":      groups,
		"apiGroupMap": apiGroupMap,
	}, c)
}

// IgnoreApi
// @Summary   忽略API
func (s *SystemApiApi) IgnoreApi(c *gin.Context) {
	var ignoreApi system.SysIgnoreApi
	err := c.ShouldBindJSON(&ignoreApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.IgnoreApi(ignoreApi)
	if err != nil {
		global.KUBEGALE_LOG.Error("忽略失败!", zap.Error(err))
		response.FailWithMessage("忽略失败", c)
		return
	}
	response.Ok(c)
}

// EnterSyncApi
// @Summary   确认同步API
func (s *SystemApiApi) EnterSyncApi(c *gin.Context) {
	var syncApi systemRes.SysSyncApis
	err := c.ShouldBindJSON(&syncApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.EnterSyncApi(syncApi)
	if err != nil {
		global.KUBEGALE_LOG.Error("忽略失败!", zap.Error(err))
		response.FailWithMessage("忽略失败", c)
		return
	}
	response.Ok(c)
}

// DeleteApi
// @Summary   删除api
func (s *SystemApiApi) DeleteApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api.KUBEGALE_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApi(api)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetApiList
// @Summary   分页获取API列表
func (s *SystemApiApi) GetApiList(c *gin.Context) {
	var pageInfo systemReq.SearchApiParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := apiService.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetApiById
// @Summary   根据id获取api
func (s *SystemApiApi) GetApiById(c *gin.Context) {
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
	api, err := apiService.GetApiById(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysAPIResponse{Api: api}, "获取成功", c)
}

// UpdateApi
// @Summary   修改基础api
func (s *SystemApiApi) UpdateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.UpdateApi(api)
	if err != nil {
		global.KUBEGALE_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetAllApis
// @Summary   获取所有的Api 不分页
func (s *SystemApiApi) GetAllApis(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	apis, err := apiService.GetAllApis(authorityID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysAPIListResponse{Apis: apis}, "获取成功", c)
}

// DeleteApisByIds
// @Summary   删除选中Api
func (s *SystemApiApi) DeleteApisByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApisByIds(ids)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FreshCasbin
// @Summary   刷新casbin缓存
func (s *SystemApiApi) FreshCasbin(c *gin.Context) {
	err := casbinService.FreshCasbin()
	if err != nil {
		global.KUBEGALE_LOG.Error("刷新失败!", zap.Error(err))
		response.FailWithMessage("刷新失败", c)
		return
	}
	response.OkWithMessage("刷新成功", c)
}
