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

type CloudLoadBalancerApi struct{}

// CloudLoadBalancerApi  同步负载均衡
func (l *CloudLoadBalancerApi) CloudLoadBalancerSync(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := cloudLoadBalancerService.SyncLoadBalancer(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("同步操作失败!", zap.Error(err))
		response.FailWithMessage("同步操作失败", c)
	} else {
		response.OkWithMessage("同步操作成功, 数据异步处理中, 请稍后!", c)
	}
}

// CloudLoadBalancerList 负载均衡列表
func (l *CloudLoadBalancerApi) CloudLoadBalancerList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchLoadBalancerParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudLoadBalancerService.List(pageInfo.LoadBalancer, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// CloudLoadBalancerTree 负载均衡Tree数据
func (l *CloudLoadBalancerApi) CloudLoadBalancerTree(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	list, err := cloudLoadBalancerService.LoadBalancerTree(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取目录树失败!", zap.Error(err))
		response.FailWithMessage("获取目录树失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}
