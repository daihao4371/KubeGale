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

type CloudVirtualMachineApi struct{}

// CloudVirtualMachineSync 同步云主机
func (p *CloudVirtualMachineApi) CloudVirtualMachineSync(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := cloudVirtualMachineService.SyncMachine(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("同步操作失败!", zap.Error(err))
		response.FailWithMessage("同步操作失败", c)
	} else {
		response.OkWithMessage("同步操作成功, 数据异步处理中, 请稍后!", c)
	}
}

// CloudVirtualMachineList 获取云主机列表
func (p *CloudVirtualMachineApi) CloudVirtualMachineList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchVirtualMachineParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudVirtualMachineService.List(pageInfo.VirtualMachine, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// CloudVirtualMachineTree
func (p *CloudVirtualMachineApi) CloudVirtualMachineTree(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	list, err := cloudVirtualMachineService.MachineTree(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取目录树失败!", zap.Error(err))
		response.FailWithMessage("获取目录树失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}
