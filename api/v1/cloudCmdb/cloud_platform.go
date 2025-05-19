package cloudCmdb

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	cloudcmdbreq "KubeGale/model/cloudCmdb/cloudcmdb"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/utils"
	cloudutils "KubeGale/utils/cloudCmdb"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudPlatformApi struct{}

// CloudPlatformList 云平台列表
func (p *CloudPlatformApi) CloudPlatformList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudPlatformService.List(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     List,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetCloudPlatformById 根据id获取CloudPlatform
func (p *CloudPlatformApi) GetCloudPlatformById(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cloud, regions, err := cloudPlatformService.GetCloudPlatformById(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(cloudcmdbreq.CloudResponse{CloudPlatform: cloud, Regions: regions}, "获取成功", c)
	}
}

// CreateCloudPlatform 创建云平台
func (p *CloudPlatformApi) CreateCloudPlatform(c *gin.Context) {
	var cloud model.CloudPlatform
	_ = c.ShouldBindJSON(&cloud)
	if err := utils.Verify(cloud, cloudutils.CloudVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := cloudPlatformService.CreateCloudPlatform(cloud); err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// UpdateCloudPlatform 更新云平台
func (p *CloudPlatformApi) UpdateCloudPlatform(c *gin.Context) {
	var cloud model.CloudPlatform
	_ = c.ShouldBindJSON(&cloud)
	if err := cloudPlatformService.UpdateCloudPlatform(cloud); err != nil {
		global.KUBEGALE_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// DeleteCloudPlatform 删除云平台
func (p *CloudPlatformApi) DeleteCloudPlatform(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := cloudPlatformService.DeleteCloudPlatform(idInfo); err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteCloudPlatformByIds 批量删除云平台
func (p *CloudPlatformApi) DeleteCloudPlatformByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := cloudPlatformService.DeleteCloudPlatformByIds(ids); err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
