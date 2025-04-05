package system

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	systemReq "KubeGale/model/system/request"
	"strconv"
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
	
	// 从上下文中获取当前用户ID
	userId, exists := c.Get("user_id")
	if exists {
		// 如果存在用户ID，则设置到操作记录中
		if userIdInt, ok := userId.(int); ok {
			sysOperationRecord.UserID = userIdInt
		} else if userIdUint, ok := userId.(uint); ok {
			sysOperationRecord.UserID = int(userIdUint)
		} else {
			// 尝试其他可能的类型转换
			global.KUBEGALE_LOG.Warn("无法将用户ID转换为整数类型")
		}
	} else {
		global.KUBEGALE_LOG.Warn("创建操作记录时未找到用户ID")
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
	// 从URL路径中获取ID参数
	idStr := c.Param("id")
	if idStr == "" {
		response.FailWithMessage("ID不能为空", c)
		return
	}

	// 处理特殊情况，如果ID以冒号开头，则去除冒号
	if len(idStr) > 0 && idStr[0] == ':' {
		idStr = idStr[1:]
	}

	// 转换ID为uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage("ID格式错误: "+err.Error(), c)
		return
	}

	// 查询记录
	reSysOperationRecord, err := operationRecordService.GetSysOperationRecord(uint(id))
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
	// 尝试从查询参数绑定
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		// 如果查询参数绑定失败，尝试从JSON绑定
		if err := c.ShouldBindJSON(&pageInfo); err != nil {
			// 设置默认值
			if pageInfo.Page <= 0 {
				pageInfo.Page = 1
			}
			if pageInfo.PageSize <= 0 {
				pageInfo.PageSize = 10
			}
		}
	}
	
	list, total, err := operationRecordService.GetSysOperationRecordInfoList(pageInfo)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	
	// 构建响应数据
	responseData := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	
	response.OkWithDetailed(responseData, "获取成功", c)
}
