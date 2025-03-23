package system

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	systemReq "KubeGale/model/system/request"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OperationRecordApi struct{}

// 创建SysOperationRecord
func (s *OperationRecordApi) CreateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.CreateSysOperationRecord(sysOperationRecord)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// 删除SysOperationRecord
func (s *OperationRecordApi) DeleteSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.DeleteSysOperationRecord(sysOperationRecord)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// 批量删除SysOperationRecord
func (s *OperationRecordApi) DeleteSysOperationRecordByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.DeleteSysOperationRecordByIds(IDS)
	if err != nil {
		global.KUBEGALE_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// 用id查询SysOperationRecord
func (s *OperationRecordApi) FindSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindQuery(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(sysOperationRecord, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reSysOperationRecord, err := operationRecordService.GetSysOperationRecord(uint(sysOperationRecord.ID))
	if err != nil {
		global.KUBEGALE_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"reSysOperationRecord": reSysOperationRecord}, "查询成功", c)
}

// 分页获取SysOperationRecord列表
func (s *OperationRecordApi) GetSysOperationRecordList(c *gin.Context) {
	var pageInfo systemReq.SysOperationRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := operationRecordService.GetSysOperationRecordInfoList(pageInfo)
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
