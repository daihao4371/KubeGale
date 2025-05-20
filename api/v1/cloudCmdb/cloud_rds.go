package cloudCmdb

import (
	"KubeGale/global"
	cloudcmdbreq "KubeGale/model/cloudCmdb/cloudcmdb"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudRDSApi struct{}

// CloudRDSSync 同步云数据库RDS信息
func (r *CloudRDSApi) CloudRDSSync(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := cloudRDSService.SyncRDS(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("同步操作失败!", zap.Error(err))
		response.FailWithMessage("同步操作失败", c)
	} else {
		response.OkWithMessage("同步操作成功, 数据异步处理中, 请稍后!", c)
	}
}

// CloudRDSList 云数据库RDS列表
func (r *CloudRDSApi) CloudRDSList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchRDSParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudRDSService.List(pageInfo.RDS, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// CloudRDSTree
func (r *CloudRDSApi) CloudRDSTree(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	list, err := cloudRDSService.RDSTree(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取目录树失败!", zap.Error(err))
		response.FailWithMessage("获取目录树失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}

// GetRDSInstance 获取RDS实例详情
func (r *CloudRDSApi) GetRDSInstance(c *gin.Context) {
	var req struct {
		Name       string `json:"name"`
		InstanceId string `json:"instanceId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 验证至少提供了一个查询参数
	if req.Name == "" && req.InstanceId == "" {
		response.FailWithMessage("请提供实例名称或实例ID", c)
		return
	}

	instance, err := cloudRDSService.GetRDSInstance(req.Name, req.InstanceId)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取RDS实例失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(instance, "获取成功", c)
}
