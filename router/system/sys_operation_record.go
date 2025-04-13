package system

import "github.com/gin-gonic/gin"

type OperationRecordRouter struct{}

func (s *OperationRecordRouter) InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("sysOperationRecord")
	{
		operationRecordRouter.POST("createSysOperationRecord", operationRecordApi.CreateSysOperationRecord)             // 新建SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecord", operationRecordApi.DeleteSysOperationRecord)           // 删除SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecordByIds", operationRecordApi.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		// 在路由配置文件中添加或修改路由
		operationRecordRouter.GET("findSysOperationRecord/:id", operationRecordApi.FindSysOperationRecord)   // 通过ID查询操作记录
		operationRecordRouter.GET("getSysOperationRecordList", operationRecordApi.GetSysOperationRecordList) // 获取SysOperationRecord列表

	}
}
