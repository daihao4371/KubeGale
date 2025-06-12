package record

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/record"
	"KubeGale/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sRecordApi struct{}

var k8sRecordService = service.ServiceGroupApp.RecordServiceGroup.K8sRecordService

func (k *K8sRecordApi) GetRecordList(c *gin.Context) {
	req := record.GetRecordListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sRecordService.GetRecordList(req); err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(record.RecordListResponse{
			Items: []interface{}{list},
			Total: total,
			PageInfo: request.PageInfo{
				Page:     req.Page,
				PageSize: req.PageSize,
				Keyword:  req.Keyword,
			},
		}, "获取成功", c)
	}
}

func (k *K8sRecordApi) DescribeRecordInfo(c *gin.Context) {
	req := record.DescribeRecordReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sRecordService.DescribeRecord(req); err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(record.DescribeRecordResponse{Items: list}, "获取成功", c)
	}
}

func (k *K8sRecordApi) UpdateRecord(c *gin.Context) {
	req := record.UpdateRecordReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sRecordService.UpdateRecord(req); err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

func (k *K8sRecordApi) DeleteRecord(c *gin.Context) {
	req := record.DeleteRecordReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sRecordService.DeleteRecord(req); err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

func (k *K8sRecordApi) CreateRecord(c *gin.Context) {
	req := record.CreateRecordReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if result, err := k8sRecordService.CreateRecord(req); err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(result, c)
	}
}
